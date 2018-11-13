import boto3
import botocore
import json
import boto
s3_client = boto3.client(
    "s3",
    aws_access_key_id="",
    aws_secret_access_key=""
)
s3 = boto.connect_s3(aws_access_key_id="",
    aws_secret_access_key="", host="s3.ap-southeast-2.amazonaws.com")
response = s3_client.list_buckets()
for bucket in response["Buckets"]:
    # Only removes the buckets with the name you want.
    if not "elasticbeanstalk-us-west-2-343345052638" == bucket["Name"] and not "terminusinflusddata" == bucket["Name"] and not "terminusinfluxdatamain"==bucket["Name"]:
        s3_objects = s3_client.list_objects_v2(Bucket=bucket["Name"])
        # Deletes the objects in the bucket before deleting the bucket.
        if "Contents" in s3_objects:
            for s3_obj in s3_objects["Contents"]:
                rm_obj = s3_client.delete_object(
                    Bucket=bucket["Name"], Key=s3_obj["Key"])
                print(rm_obj)
        buckets3 = s3.get_bucket(bucket["Name"])
        print(buckets3)
        chunk_counter = 0 #this is simply a nice to have
        keys = []
        for key in buckets3.list_versions():
            keys.append(key)
            if len(keys) > 1000:
                buckets3.delete_keys(keys)
                chunk_counter += 1
                keys = []
                print("Another 1000 done.... {n} chunks so far".format(n=chunk_counter))
        buckets3.delete_keys(keys)
        buckets3.delete()
        #rm_bucket = s3_client.delete_bucket(Bucket=bucket["Name"])
        #print(rm_bucket)
