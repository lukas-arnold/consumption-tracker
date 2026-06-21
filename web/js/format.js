const formatDate = date =>
    new Date(date).toLocaleDateString("default", {
        year: "numeric",
        month: "2-digit",
        day: "2-digit"
    });

const formatYear = date =>
    new Date(date).toLocaleDateString("default", {
        year: "numeric",
    });