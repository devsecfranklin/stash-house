# What Is the AWS Free Tier?

The AWS Free Tier is a promotion that gives free access to select AWS services for a limited time (and in
some cases, forever). It’s perfect for students, hobbyists, or anyone who wants to test the waters of
cloud computing without committing financially.

## Three types of Free Tiers

- 12-Month Free Tier: Starts when you create your AWS account.
- Always Free: Services that are permanently free within usage limits.
- Trials: Short-term, limited access to premium services.

### always free

With monthly usage caps, AWS still provides more than thirty always-free services. AWS Lambda
(1M invocations/month, 400K GB‑seconds), Amazon DynamoDB (25 GB storage + provisioned RCU/WCU),
Amazon S3 (5 GB Standard storage), Amazon CloudFront (1 TB data out + 10M requests), and Amazon
SNS (1 million publishes) are among the core, always-free services. Although the Free Plan duration
is capped (see below), these always-free limits are applicable indefinitely (not just for a full year).

## Popular Services in the AWS Free Tier

Here are some of the most useful services available at no cost (if used within the limits):
1. Amazon EC2 (Virtual Servers)
  - 750 hours/month of a t2.micro or t3.micro instance (Linux or Windows)
  - Enough to run a small server 24/7 all month long
2. Amazon S3 (Storage)
  - 5 GB of standard storage
  - 20,000 GET and 2,000 PUT requests/month
  - Great for hosting static files, images, or backups
3. Amazon RDS (Relational Database)
  - 750 hours/month of db.t2.micro (MySQL, PostgreSQL, MariaDB)
  - Ideal for experimenting with cloud databases

4. AWS Lambda (Serverless Functions)

- 1 million free requests/month
- 400,000 GB-seconds of compute time
- You can build entire backend services without a server

5. Amazon CloudFront (CDN)

- 50 GB data transfer out and 2 million HTTP/HTTPS requests per month
- Perfect for speeding up your website or app

6. Amazon DynamoDB (NoSQL Database)

- 25 GB storage and up to 200 million requests/month
- Lightweight and scalable for mobile/web apps

## What Can You Build for Free?

Here are a few real-world projects you can do entirely within the Free Tier:

- Host a static website with S3 + CloudFront
- Build a serverless API using API Gateway + Lambda + DynamoDB
- Create a small WordPress site with EC2 + RDS
- Analyze data using AWS Glue, Athena (limited usage), and S3
- Experiment with machine learning using SageMaker Studio Lab (outside Free Tier but still free)

##  Tips to Stay Within the Free Tier

- Use the AWS Billing Dashboard to track your usage.
- Set up billing alerts so you don’t accidentally go over the limits.
- Delete unused resources like EC2 instances and volumes when you’re done.
- Understand regional differences — Free Tier limits apply per account, not per region.