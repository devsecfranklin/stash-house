#ifndef lint
static char sccsid[] = "@(#)ping.c 1.17 88/02/08 SMI"; /* from UCB 4.5 4/14/86 */
#endif

/*******************************************************************************************
 *			P I N G . C
 *
 *	Compiler on LINUX
 *		cc -DLINUX ping+.c 
 *
 *	Compiler on FREE BSD
 *		cc ping+.c
 *		
 *	Compiler on Solaris 2.5: 
 *		cc ping+.c -lsocket -lnsl -lposix4
 *
 *	Usage:
 *		ping+ <host> <data size> <s/r window> <timer>
 *		s/r windown: (sender - receiver) window (similar to the TCP window concept)
 *              => it must be run as root
 *
 * Using the InterNet Control Message Protocol (ICMP) "ECHO" facility,
 * measure round-trip-delays and packet loss across network paths.
 *
 * Author -
 *	Mike Muuss
 *	U. S. Army Ballistic Research Laboratory
 *	December, 1983
 * Modified at Uc Berkeley
 * Modified at Sun Microsystems
 *
 * Paul C. Liu,  March 1990 
 * James T. Yu,  May 1999 
 * James T. Yu,  February 2003 (DePaul University)
 *
 * Status -
 *	Public Domain.  Distribution Unlimited.
 *
 * Bugs -
 *	More statistics could always be gathered.
 *	This program has to run SUID to ROOT to access the ICMP socket.
 *******************************************************************************************/

#include <signal.h>
#include <stdio.h>
#include <errno.h>
#include <time.h>
#ifdef LINUX
#include <sys/time.h>
#endif

#include <sys/param.h>
#include <sys/socket.h>
#include <sys/file.h>

#include <netinet/in_systm.h>
#include <netinet/in.h>
#include <netinet/ip.h>
#include <netinet/ip_icmp.h>
#include <netdb.h>

#define	MAXSRWIN	100
#define	MAXWAIT		5	/* max time to wait for response, sec. */
#define	MAXPACKET	32768	/* max packet size */
#ifndef MAXHOSTNAMELEN
#define MAXHOSTNAMELEN	64
#endif

int	verbose;
int	stats;
u_char	packet[MAXPACKET];
int	options;
extern	int errno;

int s;			/* Socket file descriptor */
struct hostent *hp;	/* Pointer to host info */

struct sockaddr whereto;/* Who to ping */
int datalen;		/* How much data */

char usage[] =
"Usage:	ping+ host [data size] [s/r window] [timer]\n";

char *hostname;
char hnamebuf[MAXHOSTNAMELEN];

int npackets = 999999;
int ntransmitted = 0;		/* sequence # for outbound packets = #sent */
int last_icmp_seq = 0;
int nreceived = 0;		/* # of packets we got back */
long byte_sent = 0;
long byte_recv = 0;
int ident;

int finish(), catcher();
int seconds = 5;
int timegoneby;
int elaptime = 5;
int srwindow = 5;
int lastsndcnt = 0;
int lastrcvcnt = 0;
int icmp_size = sizeof(struct icmp);

int timing = 1;	  	/* flag to compute round-trip delay */
#ifdef LINUX
	struct timeval tp;
#else
	struct timespec tp;
#endif

time_t tmin = 999999;	/* in usec */
time_t tmax = 0;
time_t tsum = 0;	/* sum of all times, for doing average */
time_t msg_tv_sec[ MAXSRWIN+1];
time_t msg_tv_usec[ MAXSRWIN+1];
time_t tvsub();

#include	<sys/utsname.h>
struct utsname	unstr;

getuname()
{
	uname(&unstr);
	/* fprintf(stdout, "%.9s", unstr.nodename); */
}


void HelpMenu() {
	printf("Usage:	ping+ host [packet size] [s/r window] [timer]\n");
	printf("\t host: hostname or IP address of the destination\n");
	printf("\t [packet size]: optional packet size.  Default=56\n");
	printf("\t\t On Ethernet LAN, set the packet size close to 1500 for the best throughput.\n");
	printf("\t [s/r window]: optional sliding window size => number of packets to send in a batch.  \n");
	printf("\t\t Use a small window size (1 or 2) to measure round trip delay  \n");
	printf("\t\t Use a large window size (10 or more) to measure throughput\n");
	printf("\t\t Congestion (too large window size) will cause long delay and low throughput\n");
	printf("\t [timer]: optional timer to stop the program after the timer.\n");
}
/*
 * 			M A I N
 */
main(argc, argv)
char *argv[];
{
	struct sockaddr_in from;
	char **av = argv;
	struct sockaddr_in *to = (struct sockaddr_in *) &whereto;
	int on = 1;
	int timeout = 20;
	struct protoent *proto;

	getuname();

	argc--, av++;
	while (argc > 0 && *av[0] == '-') {
		while (*++av[0]) switch (*av[0]) {
			case 'd':
				options |= SO_DEBUG;
				break;
			case 'r':
				options |= SO_DONTROUTE;
				break;
			case 'v':
				verbose++;
				break;
			case 's':
				stats++;
				break;
			case 'h':
				HelpMenu();
				exit(1);
				break;
		}
		argc--, av++;
	}
	if( argc < 1)  {
		printf(usage);
		exit(1);
	}

	bzero( (char *)&whereto, sizeof(struct sockaddr) );
	to->sin_family = AF_INET;
	to->sin_addr.s_addr = inet_addr(av[0]);
	if (to->sin_addr.s_addr != -1) {
		strcpy(hnamebuf, av[0]);
		hostname = hnamebuf;
	} else {
		hp = gethostbyname(av[0]);
		if (hp) {
			to->sin_family = hp->h_addrtype;
			bcopy(hp->h_addr, (caddr_t)&to->sin_addr, hp->h_length);
			hostname = hp->h_name;
		} else {
			printf("%s: unknown host %s\n", argv[0], av[0]);
			exit(1);
		}
	}

	if( argc >= 2 )
		datalen = atoi( av[1] );
	else
		datalen = 56;
	if (datalen > MAXPACKET) {
		fprintf(stderr, "ping: packet size too large\n");
		exit(1);
	}

	if (argc > 2)
	{
		srwindow = atoi(av[2]);
		if (srwindow > MAXSRWIN) srwindow = MAXSRWIN; /* max outstanding send window */
		if (srwindow < 1)  srwindow = 1; /* min outstanding send window */
	}

	if (argc > 3)
	{
		elaptime = atoi(av[3]);
		if (elaptime <= 0)
			elaptime = 1;
		if (elaptime > 600)
			elaptime = 600;		/* 10 minutes max, don't be a hogger */
		npackets =999999; /* make it arbitrary large */
	}

	if (argc > 4)
	{
		seconds = atoi(av[4]);
		if (seconds <= 0)
			seconds = 1;
		if (seconds > elaptime)
			seconds = elaptime;
	} else {
		seconds = elaptime;
	}

	ident = getpid() & 0xFFFF;

	if ((proto = getprotobyname("icmp")) == NULL) {
		fprintf(stderr, "icmp: unknown protocol\n");
		exit(10);
	}
	if ((s = socket(AF_INET, SOCK_RAW, proto->p_proto)) < 0) {
		perror("ping: socket");
		exit(5);
	}
	if (options & SO_DEBUG)
		setsockopt(s, SOL_SOCKET, SO_DEBUG, &on, sizeof(on));
		setsockopt(s, proto->p_proto, SO_OOBINLINE, &on, sizeof(on));
		on=MAXPACKET;	/* attempt 48K data buffer */
		setsockopt(s, SOL_SOCKET, SO_SNDBUF, &on, sizeof(on));
		setsockopt(s, SOL_SOCKET, SO_RCVBUF, &on, sizeof(on));
		setsockopt(s, proto->p_proto, SO_SNDBUF, &on, sizeof(on));
		setsockopt(s, proto->p_proto, SO_RCVBUF, &on, sizeof(on));
	if (options & SO_DONTROUTE)
		setsockopt(s, SOL_SOCKET, SO_DONTROUTE, &on, sizeof(on));

	if (stats)
		printf("PING %s: %d data bytes\n", hostname, datalen );

	setlinebuf( stdout );

	signal( SIGINT, finish );
	signal(SIGALRM, finish);
	if (seconds > 0) alarm(seconds);

	for (;;) {
		static int len = sizeof (packet);
		static int fromlen = sizeof (from);
		static int cc;

		while (ntransmitted < (last_icmp_seq+srwindow))
                        pinger();
		if ( (cc=recvfrom(s, packet, len, 0, &from, &fromlen)) < 0) {
			if( errno == EINTR ) fprintf(stderr, "Error in recvfrom errno=EINTR\n");
			perror("ping: recvfrom");
		}
		pr_pack( packet, cc, &from );
	}
	/*NOTREACHED*/
}

/*
 * 			P I N G E R
 * 
 * Compose and transmit an ICMP ECHO REQUEST packet.  The IP packet
 * will be added on by the kernel.  The ID field is our UNIX process ID,
 * and the sequence number is an ascending integer.  The first 8 bytes
 * of the data portion are used to hold a UNIX "timespec" struct in VAX
 * byte-order, to compute the round-trip time.
 */
pinger()
{
	static u_char outpack[MAXPACKET];
	register struct icmp *icp   = (struct icmp *) outpack;
	register u_char *datap = &outpack[8 + sizeof(struct icmp)];  /* sizeof(timespec) = 8 */
	int i, idx, cc;

	icp->icmp_type = ICMP_ECHO;
	icp->icmp_code = 0;
	icp->icmp_cksum = 0;
	icp->icmp_seq = ntransmitted;
	icp->icmp_id = ident;		/* ID */
	cc = datalen + icmp_size;		/* include ICMP portion */

	if (timing) {
#ifdef LINUX
		gettimeofday(&tp, (struct timezone *)0);
		idx = icp->icmp_seq % MAXSRWIN;
		msg_tv_sec[idx]  = tp.tv_sec;
		msg_tv_usec[idx] = tp.tv_usec;
#else
 		clock_gettime( CLOCK_REALTIME, &tp);
		idx = icp->icmp_seq % MAXSRWIN;
		msg_tv_sec[idx]  = tp.tv_sec;
		msg_tv_usec[idx] = tp.tv_nsec/1000;
#endif
	}
	if( ntransmitted == 0 ) 
		for( i=0; i<datalen; i++) *datap++ = 'x';

	/* Compute ICMP checksum here */
	icp->icmp_cksum = in_cksum( icp, cc );

	i = sendto( s, outpack, cc, 0, &whereto, sizeof(struct sockaddr) );

	if( i < 0 || i != cc )  {
		if( i<0 ) {
		    perror("sendto");
		    if (!stats)
			exit(1);
		}
		printf("ping: wrote %s %d chars, ret=%d\n", hostname, cc, i );
		fflush(stdout);
	}
	ntransmitted++;
	byte_sent += cc;
}

/*
 * 			P R _ T Y P E
 *
 * Convert an ICMP "type" field to a printable string.
 */
char *
pr_type( t )
register int t;
{
	static char *ttab[] = {
		"Echo Reply",
		"ICMP 1",
		"ICMP 2",
		"Dest Unreachable",
		"Source Quence",
		"Redirect",
		"ICMP 6",
		"ICMP 7",
		"Echo",
		"ICMP 9",
		"ICMP 10",
		"Time Exceeded",
		"Parameter Problem",
		"Timestamp",
		"Timestamp Reply",
		"Info Request",
		"Info Reply"
	};

	if( t < 0 || t > 16 )
		return("OUT-OF-RANGE");

	return(ttab[t]);
}

/*
 *			P R _ P A C K
 *
 * Print out the packet, if it came from us.  This logic is necessary
 * because ALL readers of the ICMP socket get a copy of ALL ICMP packets
 * which arrive ('tis only fair).  This permits multiple copies of this
 * program to be run without having intermingled output (or statistics!).
 */
pr_pack( buf, cc, from )
char *buf;
int cc;
struct sockaddr_in *from;
{
	struct ip *ip;
	register struct icmp *icp;
	register long *lp = (long *) packet;
	register int i, idx;
	int hlen, triptime;
	char *inet_ntoa();
	time_t delay;

	ip = (struct ip *) buf;
        hlen = ip->ip_hl << 2;

	if (cc < hlen + ICMP_MINLEN) {
		if (verbose)
			fprintf(stderr, "packet too short (%d bytes) from %s\n", cc,
				inet_ntoa(from->sin_addr));
		return;
	}
	cc -= hlen;
	icp = (struct icmp *)(buf + hlen);
	if( icp->icmp_type != ICMP_ECHOREPLY )  {
		if (verbose) {
			fprintf(stderr, "%d bytes from %s: ", cc, inet_ntoa(from->sin_addr));
			printf("icmp_type=%d (%s)\n", icp->icmp_type, pr_type(icp->icmp_type) );
			for( i=0; i<12; i++)
			    fprintf(stderr, "x%2.2x: x%8.8x\n", i*sizeof(long), *lp++ );
			fprintf(stderr, "icmp_code=%d\n", icp->icmp_code );
		}
		return;
	}
	if( icp->icmp_id != ident )
		return;			/* 'Twas not our ECHO */

	if (timing) {
#ifdef LINUX
		gettimeofday(&tp, (struct timezone *)0);
		idx = icp->icmp_seq % MAXSRWIN;
		delay = tvsub( tp.tv_sec, tp.tv_usec, msg_tv_sec[idx], msg_tv_usec[idx]);
#else
 		clock_gettime( CLOCK_REALTIME, &tp);
		idx = icp->icmp_seq % MAXSRWIN;
		delay = tvsub( tp.tv_sec, tp.tv_nsec/1000, msg_tv_sec[idx], msg_tv_usec[idx]);
#endif
		if( delay > tmax ) tmax = delay;
		if( delay < tmin ) tmin = delay;
		tsum += delay;
	}
	if (verbose)
	{
		fprintf(stderr, "Trace: %d bytes from %s w/ ", cc, inet_ntoa(from->sin_addr));
		if( timing ) fprintf(stderr, "delay = %d ", delay);
		fprintf(stderr, "icmp_seq=%d.\n", icp->icmp_seq );
	}
	if (icp->icmp_seq == 0xffff)
		last_icmp_seq = (last_icmp_seq & 0xffff0000) + 0x10000;
	else
		last_icmp_seq = (last_icmp_seq & 0xffff0000) + icp->icmp_seq + 1;

	nreceived++;
	byte_recv += (cc+hlen);
}


/*
 *			I N _ C K S U M
 *
 * Checksum routine for Internet Protocol family headers (C Version)
 *
 */
static u_short answer = 0;
in_cksum(addr, len)
u_short *addr;
int len;
{

if (answer && (ntransmitted & 0xffff)) return (--answer);
{
	register int nleft = len;
	register u_short *w = addr;
	u_short odd_byte = 0;
	register int sum = 0;
	/*
	 *  Our algorithm is simple, using a 32 bit accumulator (sum),
	 *  we add sequential 16 bit words to it, and at the end, fold
	 *  back all the carry bits from the top 16 bits into the lower
	 *  16 bits.
	 */
	while( nleft > 1 )  {
		sum += *w++;
		nleft -= 2;
	}

	/* mop up an odd byte, if necessary */
	if( nleft == 1 ) {
		*(u_char *)(&odd_byte) = *(u_char *)w;
		sum += odd_byte;
	}

	/*
	 * add back carry outs from top 16 bits to low 16 bits
	 */
	sum = (sum >> 16) + (sum & 0xffff);	/* add hi 16 to low 16 */
	sum += (sum >> 16);			/* add carry */
	answer = ~sum;				/* truncate to 16 bits */
	return (answer);
}
}

/*
 * 			T V S U B
 * 
 * Subtract 2 timeval structs:  out = out - in.
 * 
 * Out is assumed to be >= in.
 */
time_t tvsub( out_sec, out_usec, in_sec, in_usec)
time_t out_sec, out_usec, in_sec, in_usec;
{
	return( (out_usec-in_usec) + 1000000*(out_sec-in_sec) );
}

/*
 *			F I N I S H
 *
 * Print out statistics, and give up.
 * Heavily buffered STDIO is used here, so that all the statistics
 * will be written with 1 sys-write call.  This is nice when more
 * than one copy of the program is running on a terminal;  it prevents
 * the statistics output from becomming intermingled.
 */
finish()
{
	int nsnd, nrcv;
	float pass;
	double x;

	nsnd =  ntransmitted-lastsndcnt;
	nrcv =  nreceived-lastrcvcnt;
	lastsndcnt = ntransmitted;
	lastrcvcnt = nreceived;
	printf("\n----%s<===>%s PING Statistics---- time=%d,%d\n",
		unstr.nodename, hostname, timegoneby, timegoneby+seconds );
	printf("%d packets transmitted, ", nsnd );
	printf("%d packets received, ", nrcv );
	printf("%d packets lost\n", nsnd-nrcv );

	if( nsnd ) {
		pass = 100.0 - (nsnd - nrcv - srwindow)*100.0/nsnd;
		if( pass > 100.0 ) pass = 100.0;
		printf("Network Quality => %7.3f%% ICMP packets returned\n", pass);
	}
	if (nreceived && timing)
	    printf("round-trip (us)  min/avg/max = %d/%d/%d\n",
		tmin,
		tsum / nreceived,
		tmax );
	if (seconds > 0) 
		x = 1.0 * (byte_sent+byte_recv)*8.0/(1000.0*seconds);
                printf("byte sent=%d\t  byte recv=%d\t throughput (Kbps) = %8.1f\n",
			byte_sent, byte_recv, x);
	fflush(stdout);
	timegoneby += seconds;
	if (timegoneby < elaptime)
	{
		if ((elaptime - timegoneby) < seconds)
			seconds = elaptime - timegoneby;
		signal(SIGALRM, finish);
		alarm(seconds);
		if (nsnd > 0 && nrcv > 0)
			return;
	}
	printf("\nPING+ RESULT:\t%d\t%d\t%d\t%d\t%d\n",
			elaptime,
                        (byte_recv/nreceived),
			srwindow,
			ntransmitted,
			nreceived);
	exit(0);
}
