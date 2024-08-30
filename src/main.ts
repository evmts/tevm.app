import './style.css'

const querySelectorOrThrow = <T extends HTMLElement = HTMLElement>(selector: string): T => {
  const element = document.querySelector(selector);
  if (!element) {
    throw new Error(`Element not found: ${selector}`);
  }
  return element as T;
}

const mainPage = querySelectorOrThrow('#mainPage');
const applicationPage = querySelectorOrThrow('#applicationPage');
const successPage = querySelectorOrThrow('#successPage');
const showFormButton = querySelectorOrThrow<HTMLButtonElement>('#showFormButton');
const backButton = querySelectorOrThrow<HTMLButtonElement>('#backButton');
const applicationForm = querySelectorOrThrow('#applicationForm');
const backToMainButton = querySelectorOrThrow<HTMLButtonElement>('#backToMainButton');
const submitButton = querySelectorOrThrow<HTMLButtonElement>('button[type="submit"]');
const emailInput = querySelectorOrThrow<HTMLInputElement>('#emailInput');
const appsTextarea = querySelectorOrThrow<HTMLTextAreaElement>('#appsTextarea');
const developerCheckbox = querySelectorOrThrow<HTMLInputElement>('#developer');

const pages = [mainPage, applicationPage, successPage] as const

type Page = typeof pages[number];

function showPage(pageToShow: Page) {
  pages.forEach(page => {
    page.classList.remove('show');
  });
  pageToShow.classList.add('show');
  (backButton as HTMLButtonElement).style.display = pageToShow === applicationPage ? 'block' : 'none';
}

showFormButton.addEventListener('click', () => showPage(applicationPage));

backButton.addEventListener('click', (e) => {
  e.preventDefault();
  showPage(mainPage);
});

applicationForm.addEventListener('submit', function(e) {
  e.preventDefault();
  submitButton.disabled = true;
  submitButton.textContent = 'Submitting...';

  // Prepare the form data
  const formData = {
    email: emailInput.value,
    apps: appsTextarea.value,
    isDeveloper: developerCheckbox.checked
  };

  // Send the form data to the server
  fetch('/api/index', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(formData)
  })
    .then(response => {
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      return response.text(); // or response.json() if your server returns JSON
    })
    .then(data => {
      console.log('Success:', data);
      showPage(successPage);
    })
    .catch((error) => {
      console.error('Error:', error);
      alert('There was an error submitting your application. Please try again.');
      submitButton.disabled = false;
      submitButton.textContent = 'Submit Application';
    });
});

backToMainButton.addEventListener('click', () => showPage(mainPage));
