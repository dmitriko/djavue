import io
import mimetypes

from django.http import HttpResponse
from django.contrib.auth.models import User
from django.conf import settings

from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework.permissions import IsAuthenticated
from rest_framework.authtoken.models import Token

from djavue.core.models import JOB_KIND, Job, Image
from djavue.core import image_tools


def ok_response(job):
    return Response({'ok': True, 'job_id': job.id})


def error_response(error, status=400):
    return Response({'ok': False, 'error': error}, status=status)


def handle_original(job, uploaded_file):
    "Handle original job - save image as is"
    img = Image.objects.create(job=job, img=uploaded_file, kind=job.kind)
    img.save()
    return ok_response(job)


def handle_square_original(job, uploaded_file):
    '''Square of original size should not stretch the image;
    you can add white background for smaller sides

    '''
    image_db_obj = Image.objects.create(job=job, img=uploaded_file, kind=job.kind)
    image_db_obj.save()
    image_tools.make_square_original(uploaded_file, image_db_obj.img.path)
    return ok_response(job)


def handle_square_small(job, uploaded_file):
    '''Small (256px x 256px; should not stretch the image;
    you can add white background for smaller sides)

    '''
    image_db_obj = Image.objects.create(job=job, img=uploaded_file, kind=job.kind)
    image_db_obj.save()
    image_tools.make_square_small(uploaded_file, image_db_obj.img.path)
    return ok_response(job)


def handle_all_three(job, uploaded_file):
    job.kind = 'original'
    resp = handle_original(job, uploaded_file)
    if resp.status_code != 200:
        return resp
    job.kind = 'square_original'
    resp = handle_square_original(job, uploaded_file)
    if resp.status_code != 200:
        return resp
    job.kind = 'square_small'
    resp = handle_square_small(job, uploaded_file)
    if resp.status_code != 200:
        return resp
    return ok_response(job)


class RegisterUser(APIView):
    def post(self, request):
        username = request.data.get('username')
        password = request.data.get('password')
        invite_code = request.data.get('code')
        if not username or not password:
            return error_response('Missing username or password.')
        if not invite_code or invite_code != settings.INVITE_CODE:
            return error_response('Invitation code wrong or missing.')
        if User.objects.filter(username=username).count() > 0:
            return error_response('Username already exists.')
        user = User.objects.create(username=username, password=password)
        token = Token.objects.create(user=user)
        return Response({'ok': True, 'token': token.key}, status=200)


class GetImage(APIView):

    permission_classes = (IsAuthenticated,)

    def get(self, request, pk):
        try:
            image = Image.objects.get(pk=pk)
        except Image.DoesNotExist:
            return Resonse({'ok': False}, status=404)
        if image.job.user != request.user:
            return Response({'ok': False}, status=403)
        return HttpResponse(image.img, content_type=mimetypes.guess_type(
                image.img.path)[0])


class ProcessImage(APIView):
    'Perform image processing'

    permission_classes = (IsAuthenticated,)

    REQUIRED_FIELDS = ('file', 'kind')

    def get(self, request, pk):
        try:
            job = Job.objects.get(id=pk)
        except Job.DoesNotExist:
            return Response({'ok': False}, status=404)
        if job.user != request.user:
            return Response({'ok': False}, status=403)
        data = {'ok': True, 'pk': job.id, 'images':[]}
        for image in Image.objects.filter(job=job):
            data['images'].append(
                    {'pk': image.id, 'kind': image.kind}
                    )
        return Response(data)

    def post(self, request):
        errors = []
        for field in self.REQUIRED_FIELDS:
            if field not in request.data:
                errors.append('Missing {}'.format(field))
        if errors:
            return Response({'error': '\n'.join(errors)}, status=400)

        job_kind = request.data.get('kind')
        if not JOB_KIND.valid_name(job_kind):
            return Response({'error': 'Wrong job kind'}, status=400)

        job = Job.objects.create(user=request.user, kind=job_kind)
        job.save()

        if job_kind == JOB_KIND.original:
            return handle_original(job, request.data['file'])
        if job_kind == JOB_KIND.square_original:
            return handle_square_original(job, request.data['file'])
        if job_kind == JOB_KIND.square_small:
            return handle_square_small(job, request.data['file'])
        if job_kind  == JOB_KIND.all_three:
            return handle_all_three(job, request.data['file'])
        return Response({'ok':False}, status=400)

