# Cache Performance Test Script
# Save as: cache_test.ps1
# Run with: .\cache_test.ps1

Write-Host "===========================================" -ForegroundColor Cyan
Write-Host "        REDIS CACHE PERFORMANCE TEST       " -ForegroundColor Cyan
Write-Host "===========================================" -ForegroundColor Cyan
Write-Host ""

# Get user ID from user input
$userId = Read-Host "Enter User ID"

# Test endpoints with dynamic user ID
$userEndpoint = "http://localhost:5000/user/getUser/$userId?page=1"
$postsEndpoint = "http://localhost:5000/posts?page=1&id=$userId"

# Number of requests to test
$requestCount = 5

Write-Host "Testing User Profile Endpoint:"
Write-Host ""

for ($i = 1; $i -le $requestCount; $i++) {
    $time = curl.exe -s -o nul -w "%{time_total}" $userEndpoint 2>$null
    Write-Host "Request $i`: $time`s"
}

Write-Host ""
Write-Host "-------------------------------------------"
Write-Host ""

Write-Host "Testing Posts Endpoint:"
Write-Host ""

for ($i = 1; $i -le $requestCount; $i++) {
    $time = curl.exe -s -o nul -w "%{time_total}" $postsEndpoint 2>$null
    Write-Host "Request $i`: $time`s"
}

Write-Host ""
Write-Host "==========================================="
Write-Host "                TEST COMPLETE               "
Write-Host "==========================================="

# Wait for user input before closing
Read-Host "Press Enter to exit"