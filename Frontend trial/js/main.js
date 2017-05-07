const CHART = document.getElementById("radarChart");
console.log(CHART);

let radarChart = new Chart(CHART, {
    type: 'radar',
    data:data = {
        labels: ["Eating", "Drinking", "Sleeping", "Designing", "Coding", "Cycling", "Running"],
        datasets: [
            {
                label: "My First dataset",
                backgroundColor: "rgba(179,181,198,0.2)",
                borderColor: "rgba(179,181,198,1)",
                pointBackgroundColor: "rgba(179,181,198,1)",
                pointBorderColor: "#fff",
                pointHoverBackgroundColor: "#fff",
                pointHoverBorderColor: "rgba(179,181,198,1)",
                data: [65, 59, 90, 81, 56, 55, 40]
            }
        ]
    },
    options: {
        scale: {
            angleLines: {
                color: 'rgba(0,255,0,0.2)'
            },
            gridLines: {
                color: 'rgba(0,255,0,0.2)'
            }
        }
    }
});

