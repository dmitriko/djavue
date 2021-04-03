from django.test import TestCase
from django.contrib.auth.models import User
from rest_framework.test import APIClient


class TestGetToken(TestCase):

    USERNAME = 'foo'
    PASSWORD = 'bar'

    def setUp(self):
        self.user = User.objects.create_user(username=self.USERNAME,
                email='foo@bar', password=self.PASSWORD)
        self.user.save()

    def test_get_token(self):
        client = APIClient()
        resp = client.post('/api/token/',
                {'username': self.USERNAME, 'password': self.PASSWORD}, format='multipart')
        self.assertEqual(resp.status_code, 200)
        self.assertEqual(resp['Content-Type'], 'application/json')
        self.assertTrue('token' in resp.json())

