from django.test import TestCase
from django.contrib.auth.models import User
from django.urls import reverse

from rest_framework.test import APIClient

from rest_framework.authtoken.models import Token


class TestGetToken(TestCase):

    USERNAME = 'foo'
    PASSWORD = 'bar'

    def setUp(self):
        self.user = User.objects.create_user(username=self.USERNAME,
                email='foo@bar', password=self.PASSWORD)
        self.user.save()

    def test_get_token(self):
        client = APIClient()
        resp = client.post(reverse('api_token'),
                {'username': self.USERNAME, 'password': self.PASSWORD}, format='multipart')
        self.assertEqual(resp.status_code, 200)
        self.assertEqual(resp['Content-Type'], 'application/json')
        self.assertTrue('token' in resp.json())


class TestJobApiAccess(TestCase):

    def setUp(self):
        user = User.objects.create_user('foo', password='bar')
        user.save()
        self.token = Token.objects.create(user=user)
        self.token.save()

    def test_no_auth_call(self):
        client = APIClient()
        resp = client.post(reverse('api_job'), {'foo':'bar'}, format='json')
        self.assertEqual(resp.status_code, 401)

    def test_wrong_auth_call(self):
        client = APIClient()
        client.credentials(HTTP_AUTHORIZATION='Token ' + 'F00')
        resp = client.post(reverse('api_job'), {'foo':'bar'}, format='json')
        self.assertEqual(resp.status_code, 401)
        self.assertEqual(resp.json()['detail'], 'Invalid token.')

    def test_missing_field(self):
        client = APIClient()
        client.credentials(HTTP_AUTHORIZATION='Token ' + self.token.key)
        resp = client.post(reverse('api_job'), {'foo':'bar'}, format='json')
        self.assertEqual(resp.status_code, 400)

    def test_wrong_kind(self):
        client = APIClient()
        client.credentials(HTTP_AUTHORIZATION='Token ' + self.token.key)
        data = {
                'file_name':'foo.png',
                'content': 'foobar',
                'kind': 'baz'
                }
        resp = client.post(reverse('api_job'), data, format='json')
        self.assertEqual(resp.status_code, 400)
