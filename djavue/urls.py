"""djavue URL Configuration

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/3.1/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
from django.contrib import admin
from django.urls import path

from rest_framework.authtoken.views import obtain_auth_token

from djavue.core.views import ProcessImage, GetImage, RegisterUser


urlpatterns = [
    path('api/user/', RegisterUser.as_view(), name='api_user'),
    path('api/token/', obtain_auth_token, name='api_token'),
    path('api/job/', ProcessImage.as_view(), name='api_job'),
    path('api/job/<int:pk>/', ProcessImage.as_view(), name='api_job_get'),
    path('api/image/<int:pk>/', GetImage.as_view(), name='api_image_get')

]
