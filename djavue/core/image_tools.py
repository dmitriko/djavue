from PIL import Image as PILImage


def put_in_square(image, size):
    '''Use PIL image and return another PIL image
    as box, add white background for smalle side

    '''
    result_img = PILImage.new('RGB', (size, size), color=(255, 255,255))
    result_img.paste(image, (0,0))
    return result_img


def make_square_original(uploaded_file, path_to_save):
    '''Square of original size should not stretch the image;
    you can add white background for smaller sides

    '''
    orig = PILImage.open(uploaded_file)
    result_img = put_in_square(orig, max(orig.width, orig.height))
    result_img.save(path_to_save)


def make_square_small(uploaded_file, path_to_save):
    '''Small (256px x 256px; should not stretch the image;
    you can add white background for smaller sides)

    '''
    orig = PILImage.open(uploaded_file)
    cropped = orig.crop((0, 0, min(orig.width, 256), min(orig.height, 256)))
    result_img = put_in_square(cropped, 256)
    result_img.save(path_to_save)


