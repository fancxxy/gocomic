# register
http --session=fancxxy post http://127.0.0.1:8080/register usernam=fancxxy email=fancxxy@gmail.com password=fancxxy

# login
http --session=fancxxy post http://127.0.0.1:8080/login email=fancxxy@gmail.com password=fancxxy

# logout
http --session=fancxxy get http://127.0.0.1:8080/logout

# get subscribed comics 
http --session=fancxxy get http://127.0.0.1:8080/api/v1/comics

# subscribe comic
http --session=fancxxy post http://127.0.0.1:8080/api/v1/comics website=腾讯动漫 comic=航海王

# unsubscribe comic
http --session=fancxxy delete http://127.0.0.1:8080/api/v1/comics/1

# get comic and chapter list
http --session=fancxxy get http://127.0.0.1:8080/api/v1/comics/1 limit==4 page==235

# update comic and chapter list
http --session=fancxxy patch http://127.0.0.1:8080/api/v1/comics/1

# get chapter and pictures
http --session=fancxxy get http://127.0.0.1:8080/api/v1/comics/1/941