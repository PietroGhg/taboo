import glob
from PIL import Image

x = 1

files = glob.glob('cards/*')


for f in files:
    im = Image.open(f)
    w, h = im.size
    
    # Cropped image of above dimension 
    # (It will not change orginal image) 
    im1 = im.crop((0, h/2, w, h))
    im1 = im1.rotate(180)

    im1.save('newcards/' + str(x) + '.png')
    x += 1
    im2 = im.crop((0,0,w,h/2))
    im2.save('newcards/' + str(x) + '.png')

    x += 1
    
    
