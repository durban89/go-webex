<!DOCTYPE html>
<html>
<head>
  <title>Login</title>
</head>
<body>
  <form action="/login" method="post">
    Username:<input type="text" name="username"/>
    Password:<input type="password" name="password"/>
    <input type="hidden" name="token" value="{{.}}">
    <input type="submit" name="Login" />
  </form>
</body>
</html>