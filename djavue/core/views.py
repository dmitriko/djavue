from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework.permissions import IsAuthenticated


from djavue.core.models import JOB_KIND, Job, Image



def ok_response():
    return Response({'ok': True})


def handle_original(job, uploaded_file):
    "Handle original job - save image as is"
    img = Image.objects.create(job=job, img=uploaded_file)
    img.save()
    return ok_response()


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
        return Response({'ok':false}, status=400)

