docker pull redis


docker run --name redis-test-instance -p 6379:6379 -d redis


go get github.com/redis/go-redis/v8



test caching posts
 curl.exe -s -o NUL -w "ttfb=%{time_starttransfer}s total=%{time_total}s"  "http://localhost:5000/posts?page=1&id=68a3feb6d22fc774c2689e2f"

test caching user profile
 curl.exe -s -o $null -w "%{time_total}s" "http://localhost:5000/user/getUser/68a3feb6d22fc774c2689e2f?page=1"


/user/getUser/:id


....
