# geeko
A cURL-Like command line tool!

Some times I send a patch requests with commain url, params, headers and so on. Use [cURL](https://github.com/curl/curl), [Httpie](https://github.com/jkbrzt/httpie), [bat](https://github.com/astaxie/bat), I have to write the save word for each request, Even if the url is lang and lang. So I wort this.

##set command
	usage:
		set [-b|-baseUrl] url	# the base url
		set [-timeout] timeout	# the timeout for each request
		set [-t|-type] [http|form|json] 	# the request body serialization
		set [-u|-user] name 	# the user name for the server
		set [-p|-pwd|-password] pwd # the password for the server
		set [-enableCookie] enableCookie # weather use cookie

##add command
	usage:
		add [-h|-header] headers # add header for each reqeust
		add [-p|-f|-param|-form] form # add form for each reqeust

##save command
	usage:
		save [-s|-schema] name # save the schema or the last request for easy use next time.
##list command
	usage:
		list [-s|-schema] [name] # list the schema or request name

##do command 
	usage:
		do name # the request name set by save command

##use command
	usage:
		use name # the schema name set by save command

##copy method
	usage:
		copy # copy the last response to system clipboard

##state command
	usage:	
		state # the state of the tool


##http method
	usage:
		[get|post] [-h|-header header] [-p|-f|-param|-form param] url

#Thanks
Thanks for open source [httplib](https://github.com/astaxie/beego/httplib), [clipboard](github.com/atotto/clipboard), [readline](github.com/chzyer/readline).
