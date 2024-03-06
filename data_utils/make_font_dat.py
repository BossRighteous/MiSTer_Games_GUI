from PIL import Image
import struct

path = "./font.png"
# Honestly the map starts at decimal 32, no need to individual map here
char_map = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_`abcdefghijklmnopqrstuvwxyz\{|\}~█"
cell_width: int = 9
cell_height: int = 9

img = Image.open(path)
r, g, b = img.split()
result = Image.merge('RGB', (b, g, r))
img_bytes = result.tobytes()
#print(bytes)
dat_path = path+".font.dat" 
with open(dat_path, "wb") as dat:
    #add headers for img_width, img_height, cell_width,cell_height
    buffer = bytearray(2 + 2 + 1 + 1)
    struct.pack_into('<H', buffer, 0, result.width)
    struct.pack_into('<H', buffer, 2, result.height)
    struct.pack_into('<B', buffer, 4, cell_width)
    struct.pack_into('<B', buffer, 5, cell_height)
    # extend with img data
    buffer.extend(img_bytes)
    dat.write(buffer)
    dat.close()
print(dat_path + " saved successfully")