My 3 goals for this repo: 
- Try VueJS for frontend development
- Use Django for a backend, just to check what is the progress, it's funny I've made my first Django app 14 years ago. 
- Try Golang instead of Django.

For my exersise I did steal some test-your-django-docker-skill application. The task itself sounded like this:

```
Make a micro web app in Python (Django + React/Vue) that will have 3 endpoints:
Sign in
Sign up
Process image (protected)

 After sign in/sign up there should be a  form with image file input,
and a dropdown field "output" with four options: 
Original
Square of original size (should not stretch the image; you can add white background for smaller sides)
Small (256px x 256px; should not stretch the image; you can add white background for smaller sides)
"All three" (will save 3 versions of the file). 

After the user selects an image and type of the output it should send base64 of the image,
process and save it locally. The form should have backend validation: accept only image types, 
and both fields are required. 

To access the form the user should have to sign in or sign up using token based authentication. 
This should make the form endpoint protected. 

Use any libraries for the task on top of Django & React/Vue. Please use docker for the project as well.
```

Well, so far, VueJS is a real pleasure. Django as well, though all over sudden it started to give me
server error in a random manner. Hmm, no good, I am looking forward to see how Golang would handle
the same task. I started with Gin and SQLx but might try to use Beego next. 

Anyway, it is a real pleasure.
