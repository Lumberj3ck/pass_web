// UI Enhancement Scripts for Password Manager

// Loading State Management
class LoadingManager {
  constructor() {
    this.activeLoaders = new Set();
  }

  showLoading(element, text = 'Loading...') {
    const loadingEl = document.createElement('div');
    loadingEl.className = 'fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50';
    loadingEl.innerHTML = `
      <div class="bg-white rounded-lg p-6 shadow-xl">
        <div class="flex items-center">
          <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600 mr-3"></div>
          <span class="text-gray-700">${text}</span>
        </div>
      </div>
    `;
    
    document.body.appendChild(loadingEl);
    this.activeLoaders.add(loadingEl);
    
    return loadingEl;
  }

  hideLoading(loadingEl) {
    if (loadingEl && loadingEl.parentNode) {
      loadingEl.parentNode.removeChild(loadingEl);
      this.activeLoaders.delete(loadingEl);
    }
  }

  showButtonLoading(button, text = 'Processing...') {
    const originalText = button.innerHTML;
    button.classList.add('opacity-75', 'cursor-not-allowed');
    button.disabled = true;
    button.innerHTML = `
      <div class="flex items-center justify-center">
        <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-current mr-2"></div>
        ${text}
      </div>
    `;
    
    return () => {
      button.classList.remove('opacity-75', 'cursor-not-allowed');
      button.disabled = false;
      button.innerHTML = originalText;
    };
  }
}

// Form Validation and Enhancement
class FormEnhancer {
  static validateField(field) {
    const value = field.value.trim();
    const isRequired = field.hasAttribute('required');
    
    if (isRequired && !value) {
      this.showFieldError(field, 'This field is required');
      return false;
    }
    
    this.clearFieldError(field);
    return true;
  }

  static showFieldError(field, message) {
    field.classList.add('border-red-500');
    field.classList.remove('border-green-500');
    
    // Remove existing error message
    const existingError = field.parentNode.querySelector('.field-error');
    if (existingError) {
      existingError.remove();
    }
    
    const errorEl = document.createElement('div');
    errorEl.className = 'text-red-500 text-xs mt-1';
    errorEl.textContent = message;
    
    field.parentNode.appendChild(errorEl);
  }

  static clearFieldError(field) {
    field.classList.remove('border-red-500');
    field.classList.add('border-green-500');
    
    const errorEl = field.parentNode.querySelector('.field-error');
    if (errorEl) {
      errorEl.remove();
    }
  }

  static enhanceForm(form) {
    const fields = form.querySelectorAll('input, textarea');
    
    fields.forEach(field => {
      field.addEventListener('blur', () => this.validateField(field));
      field.addEventListener('input', () => {
        if (field.classList.contains('border-red-500')) {
          this.validateField(field);
        }
      });
    });

    form.addEventListener('submit', (e) => {
      let isValid = true;
      fields.forEach(field => {
        if (!this.validateField(field)) {
          isValid = false;
        }
      });
      
      if (!isValid) {
        e.preventDefault();
      }
    });
  }
}

// Alert and Notification System
class NotificationManager {
  static show(message, type = 'info', duration = 5000) {
    const alertEl = document.createElement('div');
    alertEl.className = 'fixed top-4 right-4 z-50 max-w-sm';
    
    const typeClasses = {
      success: 'bg-green-50 border-green-200 text-green-800',
      error: 'bg-red-50 border-red-200 text-red-800',
      warning: 'bg-yellow-50 border-yellow-200 text-yellow-800',
      info: 'bg-blue-50 border-blue-200 text-blue-800'
    };
    
    alertEl.innerHTML = `
      <div class="${typeClasses[type]} border rounded-lg p-4 shadow-lg">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
              ${this.getIcon(type)}
            </svg>
          </div>
          <div class="ml-3">
            <p class="text-sm font-medium">${message}</p>
          </div>
          <div class="ml-auto pl-3">
            <button class="inline-flex text-gray-400 hover:text-gray-600" onclick="this.closest('.fixed').remove()">
              <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"/>
              </svg>
            </button>
          </div>
        </div>
      </div>
    `;
    
    document.body.appendChild(alertEl);
    
    if (duration > 0) {
      setTimeout(() => {
        if (alertEl.parentNode) {
          alertEl.remove();
        }
      }, duration);
    }
    
    return alertEl;
  }

  static getIcon(type) {
    const icons = {
      success: '<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>',
      error: '<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 101.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>',
      warning: '<path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>',
      info: '<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>'
    };
    return icons[type] || icons.info;
  }
}

// Utility Functions
const Utils = {
  debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
      const later = () => {
        clearTimeout(timeout);
        func(...args);
      };
      clearTimeout(timeout);
      timeout = setTimeout(later, wait);
    };
  },

  copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
      NotificationManager.show('Copied to clipboard!', 'success', 2000);
    }).catch(() => {
      NotificationManager.show('Failed to copy to clipboard', 'error');
    });
  },

  generatePassword(length = 16) {
    const charset = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()';
    let password = '';
    for (let i = 0; i < length; i++) {
      password += charset.charAt(Math.floor(Math.random() * charset.length));
    }
    return password;
  }
};

// Global instances
window.LoadingManager = new LoadingManager();
window.FormEnhancer = FormEnhancer;
window.NotificationManager = NotificationManager;
window.Utils = Utils;

// Initialize enhancements when DOM is ready
document.addEventListener('DOMContentLoaded', function() {
  // Enhance all forms
  const forms = document.querySelectorAll('form');
  forms.forEach(form => FormEnhancer.enhanceForm(form));
  
  // Add fade-in animation to main content
  const mainContent = document.querySelector('.container, .auth-container');
  if (mainContent) {
    mainContent.classList.add('opacity-0');
    setTimeout(() => {
      mainContent.classList.add('transition-opacity', 'duration-300', 'opacity-100');
    }, 100);
  }
  
  // Add hover effects to cards
  const cards = document.querySelectorAll('.password-card, .bg-white.rounded-lg');
  cards.forEach(card => {
    card.classList.add('transition-shadow', 'duration-200');
  });
});

// HTMX event handlers for loading states
document.addEventListener('htmx:beforeRequest', function(evt) {
  const element = evt.detail.elt;
  if (element.tagName === 'BUTTON' || element.tagName === 'A') {
    element.originalText = element.innerHTML;
    element.classList.add('opacity-75', 'cursor-not-allowed');
    element.disabled = true;
    element.innerHTML = `
      <div class="flex items-center justify-center">
        <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-current mr-2"></div>
        Processing...
      </div>
    `;
  }
});

document.addEventListener('htmx:afterRequest', function(evt) {
  const element = evt.detail.elt;
  if (element.tagName === 'BUTTON' || element.tagName === 'A') {
    element.classList.remove('opacity-75', 'cursor-not-allowed');
    element.disabled = false;
    if (element.originalText) {
      element.innerHTML = element.originalText;
    }
  }
});

document.addEventListener('htmx:responseError', function(evt) {
  NotificationManager.show('An error occurred. Please try again.', 'error');
});

document.addEventListener('htmx:sendError', function(evt) {
  NotificationManager.show('Failed to connect to server. Please check your connection.', 'error');
});
    });

    form.addEventListener('submit', (e) => {
      let isValid = true;
      fields.forEach(field => {
        if (!this.validateField(field)) {
          isValid = false;
        }
      });
      
      if (!isValid) {
        e.preventDefault();
      }
    });
  }
}

// Alert and Notification System
class NotificationManager {
  static show(message, type = 'info', duration = 5000) {
    const alertEl = document.createElement('div');
    alertEl.className = `alert alert-${type} slide-in`;
    alertEl.innerHTML = `
      <div class="flex items-center justify-between">
        <div class="flex items-center">
          <svg class="h-5 w-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
            ${this.getIcon(type)}
          </svg>
          <span>${message}</span>
        </div>
        <button class="ml-4 text-sm font-medium hover:underline" onclick="this.parentNode.parentNode.remove()">
          Dismiss
        </button>
      </div>
    `;
    
    document.body.appendChild(alertEl);
    
    if (duration > 0) {
      setTimeout(() => {
        if (alertEl.parentNode) {
          alertEl.remove();
        }
      }, duration);
    }
    
    return alertEl;
  }

  static getIcon(type) {
    const icons = {
      success: '<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>',
      error: '<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 101.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>',
      warning: '<path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>',
      info: '<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>'
    };
    return icons[type] || icons.info;
  }
}

// Utility Functions
const Utils = {
  debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
      const later = () => {
        clearTimeout(timeout);
        func(...args);
      };
      clearTimeout(timeout);
      timeout = setTimeout(later, wait);
    };
  },

  copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
      NotificationManager.show('Copied to clipboard!', 'success', 2000);
    }).catch(() => {
      NotificationManager.show('Failed to copy to clipboard', 'error');
    });
  },

  generatePassword(length = 16) {
    const charset = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()';
    let password = '';
    for (let i = 0; i < length; i++) {
      password += charset.charAt(Math.floor(Math.random() * charset.length));
    }
    return password;
  }
};

// Global instances
window.LoadingManager = new LoadingManager();
window.FormEnhancer = FormEnhancer;
window.NotificationManager = NotificationManager;
window.Utils = Utils;

// Initialize enhancements when DOM is ready
document.addEventListener('DOMContentLoaded', function() {
  // Enhance all forms
  const forms = document.querySelectorAll('form');
  forms.forEach(form => FormEnhancer.enhanceForm(form));
  
  // Add fade-in animation to main content
  const mainContent = document.querySelector('.container, .auth-container');
  if (mainContent) {
    mainContent.classList.add('fade-in');
  }
  
  // Add hover effects to cards
  const cards = document.querySelectorAll('.password-card, .card');
  cards.forEach(card => {
    card.classList.add('hover-lift');
  });
});

// HTMX event handlers for loading states
document.addEventListener('htmx:beforeRequest', function(evt) {
  const element = evt.detail.elt;
  if (element.tagName === 'BUTTON' || element.tagName === 'A') {
    element.originalText = element.innerHTML;
    element.classList.add('btn-loading');
    element.disabled = true;
  }
});

document.addEventListener('htmx:afterRequest', function(evt) {
  const element = evt.detail.elt;
  if (element.tagName === 'BUTTON' || element.tagName === 'A') {
    element.classList.remove('btn-loading');
    element.disabled = false;
    if (element.originalText) {
      element.innerHTML = element.originalText;
    }
  }
});

document.addEventListener('htmx:responseError', function(evt) {
  NotificationManager.show('An error occurred. Please try again.', 'error');
});

document.addEventListener('htmx:sendError', function(evt) {
  NotificationManager.show('Failed to connect to server. Please check your connection.', 'error');
});