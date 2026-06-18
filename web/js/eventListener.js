document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll(".date")
        .forEach(el => el.textContent = formatDate(el.textContent));
    const currentDate = new Date(Date.now() - new Date().getTimezoneOffset() * 60000)
        .toISOString()
        .split('T')[0];

    const lastYear = new Date().getFullYear() - 1;

    document.querySelectorAll('input[type="date"]').forEach(input => {
        if (input.value) return;

        switch (input.id) {
            case 'lastYearStart':
                input.value = `${lastYear}-01-01`;
                break;
            case 'lastYearEnd':
                input.value = `${lastYear}-12-31`;
                break;
            default:
                input.value = currentDate;
        }
    });

    const yearInput = document.getElementById('lastYear');
    if (yearInput && !yearInput.value) {
        yearInput.value = lastYear;
    }
});