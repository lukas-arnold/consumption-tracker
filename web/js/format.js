document.addEventListener('DOMContentLoaded', function() {
    formatDates();
    formatFloats();
});

function formatDates() {
    const dateElements = document.querySelectorAll(".date");
    
    dateElements.forEach(element => {
        const originalDate = new Date(element.textContent);
        const options = {year: "numeric", month: "2-digit", day: "2-digit"};
        element.textContent = originalDate.toLocaleDateString("default", options);
    });
}

function formatFloats() {
    const floatElements = document.querySelectorAll(".float-value");
    
    floatElements.forEach(element => {
        let value = parseFloat(element.textContent);
        if (!isNaN(value)) {
            element.textContent = value.toLocaleString("de-DE", {minimumFractionDigits: 2, maximumFractionDigits: 2});
        }
    });
}

function formatDate(date) {
    date = new Date(date);
    const options = {year: "numeric", month: "2-digit", day: "2-digit"};
    return date.toLocaleDateString("default", options);
}

function formatFloat(value) {
    const num = parseFloat(value);
    if (isNaN(num)) return value;
    return num.toLocaleString("de-DE", {minimumFractionDigits: 2, maximumFractionDigits: 2});
}
