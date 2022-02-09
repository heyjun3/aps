var drawChart = function (id, date, price, rank) {
    const ctx = document.getElementById(id);
    var myLineChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: date,
            datasets: [
                {
                    label: 'Price',
                    data: price,
                    // データライン
                    borderColor: 'blue',
                },
                {
                    label: 'Rank',
                    data: rank,
                    // データライン
                    borderColor: 'green',
                    yAxisID: 'y2',
                },
            ],
        },
        options: {
            y2: {
                position: 'right',
            }
        }
    });
}

var sample = function(){
    console.log('test')
}