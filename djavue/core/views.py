from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework.permissions import IsAuthenticated


class ProcessImage(APIView):
    'Perform image processing'

    permission_classes = (IsAuthenticated,)

    REQUIRED_FIELDS = ('file_name', 'content', 'kind')
    KINDS = ('original', 'square_original', 'square_small', 'all_three')

    def post(self, request):
        errors = []
        for field in self.REQUIRED_FIELDS:
            if field not in request.data:
                errors.append('Missing {}'.format(field))
        if errors:
            return Response({'error': "\n".join(errors)}, status=400)
        return Response({'foo':'bar'})

