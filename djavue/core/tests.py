import pathlib
import os

from django.test import TestCase
from django.contrib.auth.models import User
from django.urls import reverse
from django.core.files.uploadedfile import SimpleUploadedFile
from django.conf import settings

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


class TestRegisterUser(TestCase):
    def setUp(self):
        self.url = reverse('api_user')
        self.client = APIClient()

    def test_missing_username_or_password(self):
        for item in  [{'code':'1'},
                {'username':'foo', 'code': '1'},
                {'password': 'foo', 'code':'1'}]:
            resp = self.client.post(self.url, item, format='json')
            self.assertEqual(resp.status_code, 400)
            self.assertFalse(resp.json()['ok'])
            self.assertEqual(resp.json().get('error'), 'Missing username or password.')

    def test_wrong_or_missing_invite_code(self):
        for item in [{'username':'u', 'password':'p', 'code':'WRONG'},
                {'username':'u', 'password': 'p'}]:
            resp = self.client.post(self.url, {'username': 'foo', 'password': 'bar'})
            self.assertEqual(resp.status_code, 400)
            self.assertFalse(resp.json()['ok'])
            self.assertEqual(resp.json().get('error'), 'Invitation code wrong or missing.')

    def test_already_exists(self):
        foo_user = User.objects.create_user('foo', password='bar')
        resp = self.client.post(self.url, {
            'username': 'foo', 'password': 'bar', 'code':settings.INVITE_CODE})
        self.assertEqual(resp.status_code, 400)
        self.assertFalse(resp.json()['ok'])
        self.assertEqual(resp.json().get('error'), 'Username already exists.')

    def test_register_ok(self):
        resp = self.client.post(self.url, {
            'username': 'foo', 'password': 'bar', 'code':settings.INVITE_CODE})
        self.assertEqual(resp.status_code, 200)
        user = User.objects.get(username='foo')
        token, _ = Token.objects.get_or_create(user=user)
        self.assertEqual(resp.json().get('token'), token.key)


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

    def tearDown(self):
        for image in Image.objects.all():
            os.unlink(image.img.path)

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

    def test_square_original(self):
        resp = self.client.post(reverse('api_job'),
            dict(file=self.upload_file,
                kind='square_original'))
        self.assertEqual(resp.status_code, 200)
        image = Image.objects.get(job__id=resp.json()['job_id'])
        self.assertEqual(image.img.width, image.img.height)

    def test_square_small(self):
        resp = self.client.post(reverse('api_job'),
            dict(file=self.upload_file,
                kind='square_small'))
        self.assertEqual(resp.status_code, 200)
        image = Image.objects.get(job__id=resp.json()['job_id'])
        self.assertEqual(256, image.img.width)
        self.assertEqual(256, image.img.height)

    def test_all_three(self):
        resp = self.client.post(reverse('api_job'),
            dict(file=self.upload_file,
                kind='all_three'))
        self.assertEqual(resp.status_code, 200)
        image1 = Image.objects.get(job__id=resp.json()['job_id'], kind='original')
        image2 = Image.objects.get(job__id=resp.json()['job_id'], kind='square_original')
        image3 = Image.objects.get(job__id=resp.json()['job_id'], kind='square_small')
        self.assertEqual(image1.img.size, self.upload_file.size)
        self.assertEqual(image2.img.width, image2.img.height)
        self.assertEqual(image3.img.width, image3.img.height, 256)

    def test_get_all_three(self):
        "Process file and get all three images back"
        # test 404 error first
        resp = self.client.get(reverse('api_job_get', args=(111111,)))
        self.assertEqual(resp.status_code, 404)
        resp = self.client.post(reverse('api_job'),
            dict(file=self.upload_file,
                kind='all_three'))
        job_id = resp.json()['job_id']
        get_url = reverse('api_job_get', args=(job_id,))
        resp = self.client.get(get_url)
        self.assertEqual(resp.status_code, 200)
        data = resp.json()
        images = {}
        for item in data['images']:
            images[item['kind']] = item['pk']
        self.assertEqual(len(images), 3)
        for image_kind, image_pk in images.items():
            self.assertTrue(Image.objects.get(pk=image_pk, kind=image_kind))

    def test_image_get(self):
        resp = self.client.post(reverse('api_job'),
            dict(file=self.upload_file,
                kind='square_small'))
        self.assertEqual(resp.status_code, 200)
        job_id = resp.json()['job_id']
        stored_image = Image.objects.get(job__id=job_id)
        resp = self.client.get(reverse('api_image_get', args=(stored_image.id,)))
        self.assertEqual(resp.status_code, 200)
        self.assertEqual(len(resp.content), stored_image.img.size)
#        import IPython; IPython.embed(using=False)

