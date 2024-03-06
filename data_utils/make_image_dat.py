from PIL import Image
import struct

path = "./jeff.png"

img = Image.open(path)
r, g, b = img.split()
result = Image.merge('RGB', (b, g, r))
img_bytes = result.tobytes()
#print(bytes)
dat_path = path+".dat" 
with open(dat_path, "wb") as dat:
    #add headers for img_width, img_height, cell_width,cell_height
    buffer = bytearray(2 + 2)
    struct.pack_into('<H', buffer, 0, result.width)
    struct.pack_into('<H', buffer, 2, result.height)
    # extend with img data
    buffer.extend(img_bytes)
    dat.write(buffer)
    dat.close()
print(dat_path + " saved successfully")