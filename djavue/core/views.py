from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework.permissions import IsAuthenticated


from djavue.core.models import Job, JOB_KINDS


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
            return Response({'error': "\n".join(errors)}, status=400)
        if request.data['kind'] not in JOB_KINDS:
            return Response({'error': 'Wrong job kind'}, status=400)
        job = Job.objects.create(user=request.user, kind=request.data['kind'])
        job.save()
        return Response({'ok':True})

