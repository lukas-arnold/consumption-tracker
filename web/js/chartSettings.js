Chart.defaults.elements.line.tension = 0.4;
Chart.defaults.interaction = {
    mode: 'index',
    intersect: false
};

function renderChart(canvasId, model) {
    const ctx = document.getElementById(canvasId);

    const datasets = model.sets.map(s => ({
        label: s.label,
        data: s.data
    }));

    const hasMultipleSets = model.sets.length > 1;

    const options = {
        responsive: true,
        maintainAspectRatio: false,
        interaction: {
            mode: "index",
            intersect: false
        },
        plugins: {
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const ds = model.sets[context.datasetIndex];
                        const value = context.raw;

                        if (!ds.unit) return context.dataset.label + ": " + value;

                        return context.dataset.label + ": " + value + " " + ds.unit;
                    }
                }
            }
        },
        scales: {
            y: {
                beginAtZero: true
            }
        }
    };

    // If there are multiple sets, assign the second set to a secondary y-axis
    if (hasMultipleSets) {
        options.scales.y1 = {
            position: "right",
            grid: {
                drawOnChartArea: false
            }
        };
    }

    new Chart(ctx, {
        type: model.type || "line",
        data: {
            labels: model.labels,
            datasets
        },
        options
    });
}