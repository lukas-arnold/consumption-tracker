Chart.defaults.elements.line.tension = 0.4;
Chart.defaults.interaction = {
    mode: 'index',
    intersect: false
};

function isIsoDate(value) {
    return typeof value === "string" &&
        /^\d{4}-\d{2}-\d{2}$/.test(value);
}

function formatValue(value, unit) {
    if (value == null) return value;

    const num = Number(value);

    if (unit === "€" || unit === "€/kWh" || unit === "€/l" || unit === "€/m³") {
        return num.toLocaleString("de-DE", {
            minimumFractionDigits: 2,
            maximumFractionDigits: 2
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
                beginAtZero: true,
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

                            if (isIsoDate(label)) {
                                return formatDate(label);
                            }

                            return label;
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

                            if (isIsoDate(label)) {
                                return formatDate(label);
                            }

                            return label;
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