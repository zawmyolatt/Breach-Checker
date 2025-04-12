document.addEventListener('DOMContentLoaded', function() {
    // Form validation
    const form = document.querySelector('form');
    if (form) {
        form.addEventListener('submit', function(e) {
            const emailInput = document.getElementById('email');
            if (!emailInput.value.trim()) {
                e.preventDefault();
                alert('Please enter an email address');
                return;
            }
            
            const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailPattern.test(emailInput.value)) {
                e.preventDefault();
                alert('Please enter a valid email address');
                return;
            }
        });
    }
    
    // Add animation to result page
    const resultSection = document.querySelector('.result');
    if (resultSection) {
        resultSection.classList.add('fade-in');
    }
}); 