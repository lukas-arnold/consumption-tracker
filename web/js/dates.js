document.addEventListener('DOMContentLoaded', function() {
    const today = new Date();
    const offset = today.getTimezoneOffset();
    const todayUTC = new Date(today.getTime() + offset * 60 * 1000);
    const currentDate = todayUTC.toISOString().split('T')[0];

    const lastYear = new Date().getFullYear() - 1;
    const lastYearStart = `${lastYear}-01-01`;
    const lastYearEnd = `${lastYear}-12-31`;

    document.querySelectorAll('input[type="date"]').forEach(input => {
        if (!input.value) {
            input.value = currentDate;
            if (input.id === 'lastYearStart') {
                input.value = lastYearStart;
            };
            if (input.id === 'lastYearEnd') {
                input.value = lastYearEnd;
            };
        }
    });

    const yearInput = document.getElementById('lastYear');
    if (yearInput && !yearInput.value) {
        yearInput.value = lastYear;
    }
});
