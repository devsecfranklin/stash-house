
local base(o) = o + {
    db +: {
        password : "175k2v4o3w482c8ela9ww",
		host : "localhost",
		port : 54321,
    },
	docker_compose : false,
	external_addr : "www.bitsmasher.net",
	top_dir : ".",
};
local final(o) = o + {
	cks : { enc_keys : ["Sid6W0tDNzOjJcCojOuhbrL6YeMce3rqTnHNQae3QP2"] }
};
{ base : base, final : final }
	