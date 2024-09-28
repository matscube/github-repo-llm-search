

* GitHub API
  - https://api.github.com/search/repositories?q=stars:%3E100000
  - https://api.github.com/search/repositories?q=stars:%3E10000&sort=stars&per_page=100&page=1
* GitHub search
  - https://github.com/search?q=stars%3A%3E1000&type=repositories




# multi filtering conditions are NOT available (like star and created)
 https://api.github.com/search/repositories?q=stars:>1000&created:>2020-07-01&sort=updated

 https://api.github.com/search/repositories?q=created:2020-07-01..2020-07-02&sort=created




# https://api.github.com/search/repositories?q=created:%3C2025-01-01&sort=created
=> 57317167 total
# https://api.github.com/search/repositories?q=stars:%3E=0
=> 70629692 total

> 30000: total: 660
> 20000: per 10000 + 1000
> 10000: per 1000 + 100
> 6000: per 200 + 50
> 2000: per 40 + 10
> 1000: per 10 + 5
> 500: per 3 + 2

https://api.github.com/search/repositories?q=stars:>50000&sort=stars => 251
https://api.github.com/search/repositories?q=stars:30000..60000&sort=stars => 481
https://api.github.com/search/repositories?q=stars:20000..31000&sort=stars => 772
https://api.github.com/search/repositories?q=stars:19000..21000&sort=stars



https://api.github.com/search/repositories?q=stars:1000..1002&sort=stars



* get base path sliding window
* loop for base path with paging
  * store items 
