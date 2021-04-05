from django.db import models
from django.contrib.auth.models import User



class JOB_KIND:
    original = 'original'
    square_original = 'square_original'
    square_small = 'square_small'
    all_three = 'all_three'

    @classmethod
    def valid_name(cls, name):
        if not name.startswith('_'):
            return name in cls.__dict__
        return False


class Job(models.Model):
    user = models.ForeignKey(User, on_delete=models.CASCADE, null=True)
    kind = models.CharField(max_length=50, choices=[(x, x) for x in [JOB_KIND.original,
        JOB_KIND.square_original, JOB_KIND.square_small, JOB_KIND.all_three]])


class Image(models.Model):
    job = models.ForeignKey(Job, on_delete=models.CASCADE)
    img = models.ImageField(upload_to='images')
    kind = models.CharField(max_length=50, choices=[(x, x) for x in [JOB_KIND.original,
        JOB_KIND.square_original, JOB_KIND.square_small]])
