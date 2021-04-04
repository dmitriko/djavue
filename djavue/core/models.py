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

    @classmethod
    def choices(cls, include_all=True):
        all_values = [y for x,y in cls.__dict__.items() if not x.startswith('_')]
        if include_all:
            return [(x, x) for x in all_values]
        return [(x, x) for x in all_values if x != cls.all_three]


class Job(models.Model):
    user = models.ForeignKey(User, on_delete=models.CASCADE, null=True)
    kind = models.CharField(max_length=50, choices=JOB_KIND.choices())


class Image(models.Model):
    job = models.ForeignKey(Job, on_delete=models.CASCADE)
    img = models.ImageField(upload_to='images')
    kind = models.CharField(max_length=50, choices=JOB_KIND.choices(include_all=False))
