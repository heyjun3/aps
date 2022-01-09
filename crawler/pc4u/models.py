from crawler.models import Product
from crawler.models import Base
from crawler.models import postgresql_engine


class Pc4u(Product, Base):
    __tablename__ = 'pc4u_products'

    

def init_db():
    Base.metadata.create_all(bind=postgresql_engine)