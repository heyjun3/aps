import datetime

import pandas as pd
import numpy as np
from sqlalchemy.engine.default import DefaultExecutionContext


def recharts_data(context: dict|DefaultExecutionContext) -> dict|None:
    if isinstance(context, DefaultExecutionContext):
        params = context.get_current_parameters()
    elif isinstance(context, dict):
        params = context
    else:
        raise Exception

    rank_data = params.get('rank_data')
    price_data = params.get('price_data')
    if rank_data is None or price_data is None:
        return None

    today = datetime.datetime.now().date()
    start_date = today - datetime.timedelta(days=90)
    end_date = today
    date_index = pd.date_range(start_date, end_date)
    date_index_df = pd.DataFrame(data=date_index, columns=['date'])
    date_index_df = date_index_df['date'].dt.date

    price_df = pd.DataFrame(data=price_data.items(), columns=['date', 'price']).astype({'date': int, 'price': int})
    price_df['date'] = price_df['date'].map(keepa_time_to_datetime_date)
    rank_df = pd.DataFrame(data=rank_data.items(), columns=['date', 'rank']).astype({'date': int, 'rank': int})
    rank_df['date'] = rank_df['date'].map(keepa_time_to_datetime_date)

    df = pd.merge(date_index_df, price_df, on='date', how='outer')
    df = pd.merge(df, rank_df, on='date', how='outer')
    df = df.replace(-1.0, np.nan)
    df = df.sort_values('date', ascending=True)
    df = df.fillna(method='ffill')
    df = df.fillna(method='bfill')
    df = df.replace([np.nan], [None])
    df = df[df['date'] > start_date]
    df = df.sort_values('date', ascending=True)
    df['date'] = df['date'].map(lambda x: x.strftime('%Y-%m-%d'))
    data = df.to_dict(orient='records')

    return {'data': data}

def keepa_time_to_datetime_date(keepa_time: int):
    unix_time = (keepa_time + 21564000) * 60
    date_time = datetime.datetime.fromtimestamp(unix_time)
    return date_time.date()


def unix_time_to_keepa_time(unix_time: float) -> str:
    keepa_time = round(unix_time / 60 - 21564000)
    return str(keepa_time)
