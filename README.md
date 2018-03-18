# TECT
Backend for Proteus clients

TECT (c) 2018 by [Krogoth](https://twitter.com/le_krogoth) of [Ministry of Zombie Defense](http://www.mzd.org.uk/) and contributing authors

# Build
In root folder:

> gb build

Builds binary into ~/bin folder.

# Run

> cd bin
> ./TECT

Make sure that you have the following structure:

> /bin/  
> ---/TECT -> executable  
> ---/TECT.config.json -> will be generated if not existent  
> ---/TECT.db -> sqlite3, automatically generated
> ---/firmware/  -> your firmware files here  
> ---/spiffs/  -> your spiffs files here  

TECT will always verify the latest (as in lexical ordering) file in the firmware or spiffs folder and compare the md5 given in the GET request against the md5 of the latest file on disc.

If the md5 matches, TECT will return a not changed status code. If the md5 does not match, TECT will always return the file without checking if the given md5 represents an earlier file or not. 

Be aware that if you mess up the firmware folder, you could theoretically install an older file since TECT does not check the version history.

# Prerequisites
To build TECT, you will need a somewhat current Golang installation as well as GB installed.

# Licence
TECT is published under the AGPLv3. See LICENCE file for details. 