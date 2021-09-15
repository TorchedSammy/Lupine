local lupinelog = require 'lupinelog'

os.execute 'hostname -F /etc/hostname'
lupinelog.success 'Set system hostname'
