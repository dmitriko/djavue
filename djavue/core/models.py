from django.db import models
from django.contrib.auth.models import User


JOB_KINDS = ('original', 'square_original', 'square_small', 'all_three')


class Job(models.Model):
    user = models.ForeignKey(User, on_delete=models.CASCADE, null=True)
    kind = models.CharField(max_length=50, choices=[(x, x) for x in JOB_KINDS])

