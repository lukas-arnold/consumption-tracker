document.addEventListener('DOMContentLoaded', function() {
    document.querySelectorAll(".date")
        .forEach(el => el.textContent = formatDate(el.textContent));
});

const formatDate = date =>
    new Date(date).toLocaleDateString("default", {
        year: "numeric",
        month: "2-digit",
        day: "2-digit"
    });