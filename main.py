import sys

from mws import api
from keepa import keepa
from crawler.buffalo import buffalo
from crawler.pc4u import pc4u
from crawler.rakuten import rakuten
from crawler.super import super
from crawler.netsea import netsea
from ims import repeat
from ims import monthly


if __name__ == '__main__':

    args = sys.argv

    if args[1] == 'keepa':
        keepa.main()
    elif args[1] == 'mws':
        api.main()
    elif args[1] == 'buffalo':
        buffalo.main()
    elif args[1] == 'pc4u':
        pc4u.schedule_pc4u_task_everyday()
    elif args[1] == 'rakuten':
        rakuten.schedule()
    elif args[1] == 'super':
        super.super_all()
    elif args[1] == 'netsea':
        netsea.netsea_all()
    elif args[1] == 'repeat':
        repeat.main()
    elif args[1] == 'monthly':
        monthly.main()
    else:
        sys.stdout.write(f'{args[1]} is not a command')
