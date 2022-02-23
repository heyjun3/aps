from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from sqlalchemy.ext.declarative import declarative_base


engine = create_engine('sqlite:///')
Session = sessionmaker(bine=engine)
Base = declarative_base()

Base.metadata.create_all(bind=engine)



