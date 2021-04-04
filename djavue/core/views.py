import io

from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework.permissions import IsAuthenticated

from PIL import Image as PILImage

from djavue.core.models import JOB_KIND, Job, Image



def ok_response(job):
    return Response({'ok': True, 'job_id': job.id})


def handle_original(job, uploaded_file):
    "Handle original job - save image as is"
    img = Image.objects.create(job=job, img=uploaded_file, kind=job.kind)
    img.save()
    return ok_response(job)


def put_in_square(image, size):
    '''Use PIL image and return another PIL image
    as box, add white background for smalle side

    '''
    result_img = PILImage.new('RGB', (size, size), color=(255, 255,255))
    result_img.paste(image, (0,0))
    return result_img


def get_pil_format(uploaded_file):
    "Return JPEG or PNG"
    if 'png' in uploaded_file.content_type:
        return 'PNG'
    return 'JPEG'


def handle_square_original(job, uploaded_file):
    '''Square of original size should not stretch the image;
    you can add white background for smaller sides

    '''
    image_db_obj = Image.objects.create(job=job, img=uploaded_file, kind=job.kind)
    image_db_obj.save()
    orig = PILImage.open(uploaded_file)
    result_img = put_in_square(orig, max(orig.width, orig.height))
    result_img.save(image_db_obj.img.path)
    return ok_response(job)


def handle_square_small(job, uploaded_file):
    '''Small (256px x 256px; should not stretch the image;
    you can add white background for smaller sides)

    '''
    image_db_obj = Image.objects.create(job=job, img=uploaded_file, kind=job.kind)
    image_db_obj.save()
    orig = PILImage.open(uploaded_file)
    cropped = orig.crop((0, 0, min(orig.width, 256), min(orig.height, 256)))
    result_img = put_in_square(cropped, 256)
    result_img.save(image_db_obj.img.path)
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



class ProcessImage(APIView):
    'Perform image processing'

    permission_classes = (IsAuthenticated,)

    REQUIRED_FIELDS = ('file', 'kind')

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

