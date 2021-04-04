import base64
import pathlib

from django.test import TestCase
from django.contrib.auth.models import User
from django.urls import reverse
from django.core.files.uploadedfile import SimpleUploadedFile

from rest_framework.test import APIClient

from rest_framework.authtoken.models import Token

from djavue.core.models import Job, Image


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
                'file': SimpleUploadedFile(
                    'foo.txt', b'foo', content_type='text/plain'),
                'kind': 'baz'
                }
        resp = client.post(reverse('api_job'), data, format='json')
        self.assertEqual(resp.status_code, 400)


class TestJobApiProcess(TestCase):

    def setUp(self):
        user = User.objects.create_user('foo', password='bar')
        user.save()
        token = Token.objects.create(user=user)
        token.save()
        self.client = APIClient()
        self.client.credentials(HTTP_AUTHORIZATION='Token ' + token.key)
        path = pathlib.Path(__file__).parent / 'test_data' / 'img.png'
        with open(path, 'rb') as f:
            self.upload_file = SimpleUploadedFile(
                    'foo.png', f.read(), content_type='image/png')

    def test_original(self):
        resp = self.client.post(reverse('api_job'),
            dict(file=self.upload_file,
                kind='original'))
        self.assertEqual(resp.status_code, 200)
        self.assertEqual(Job.objects.count(), 1)
        job = list(Job.objects.all())[0]
        self.assertTrue(job.user is not None)
        self.assertEqual(job.user.username, 'foo')
        self.assertEqual(job.kind, 'original')
        self.assertEqual(Image.objects.filter(job=job).count(), 1)
        image = Image.objects.get(job=job)
        self.assertEqual(image.img.size, self.upload_file.size)

