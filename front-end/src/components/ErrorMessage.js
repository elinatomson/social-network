export function displayErrorMessage(message) {
    var formContainer = document.getElementById('error');
    var errorContainer = document.createElement('div');
    errorContainer.className = 'message';
    errorContainer.textContent = message;
    formContainer.appendChild(errorContainer);
}