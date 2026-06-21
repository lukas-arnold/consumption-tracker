Chart.defaults.elements.line.tension = 0.4;
Chart.defaults.interaction = {
    mode: 'index',
    intersect: false
};

function formatValue(value, unit) {
    if (value == null) return value;

    const num = Number(value);

    if (unit === "€") {
        return num.toLocaleString("de-DE", {
            minimumFractionDigits: 2,
            maximumFractionDigits: 2
        }) + " " + unit;
    }

    if (unit === "€/kWh" || unit === "€/l" || unit === "€/m³") {
        return num.toLocaleString("de-DE", {
            minimumFractionDigits: 3,
            maximumFractionDigits: 3
        }) + " " + unit;
    }

    if (unit === "kWh" || unit === "l" || unit === "m³") {
        return num.toLocaleString("de-DE", {
            maximumFractionDigits: 0
        }) + " " + unit;
    }

    if (unit === "cm") {
        return num.toLocaleString("de-DE", {
            maximumFractionDigits: 1
        }) + " cm";
    }

    if (unit === "%") {
        return num.toLocaleString("de-DE", {
            maximumFractionDigits: 2
        }) + " %";
    }

    return num;
}

function renderChart(canvasId, model) {
    const ctx = document.getElementById(canvasId);

    const datasets = model.sets.map(s => ({
        label: s.label,
        data: s.data,
        yAxisID: s.yAxis || "y",
        hidden: s.hidden || false
    }));

    const yAxes = {};

    model.sets.forEach(s => {
        const axis = s.yAxis || "y";
        if (!yAxes[axis]) {
            yAxes[axis] = {
                type: "linear",
                position: axis === "y1" ? "right" : "left",

                // Fix y axis only for fillChart
                ...(canvasId === "fillChart" && axis === "y" ? {
                    min: 0,
                    max: 150
                } : {}),
                ...(canvasId === "fillChart" && axis === "y1" ? {
                    min: 0,
                    max: 100
                } : {}),

                grid: axis === "y1" ? { drawOnChartArea: false } : undefined,
                ticks: {
                    callback: function(value) {
                        return formatValue(value, s.unit);
                    }
                }
            };
        }
    });

    new Chart(ctx, {
        type: model.type || "line",
        data: {
            labels: model.labels,
            datasets
        },
        options: {
            scales: {
                x: {
                    ticks: {
                        callback: function(value) {
                            const label = this.getLabelForValue(value);

                            return canvasId === "fillChart"
                                ? formatDate(label)
                                : formatYear(label);
                        }
                    }
                },
                ...yAxes
            },
            plugins: {
                tooltip: {
                    callbacks: {
                        title: (items) => {
                            const label = items[0].label;

                            return canvasId === "fillChart"
                                ? formatDate(label)
                                : formatYear(label);
                        },
                        label: (ctx) => {
                            const unit = model.sets[ctx.datasetIndex]?.unit || "";
                            return `${ctx.dataset.label}: ${formatValue(ctx.raw, unit)}`;
                        }
                    }
                }
            }
        }
    });
}