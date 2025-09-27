const API_URL = ''; // Same origin

function parseJwt(token) {
    try {
        return JSON.parse(atob(token.split('.')[1]));
    } catch (e) {
        return null;
    }
}

async function login() {
    const username = document.getElementById('username').value;
    const errorMessage = document.getElementById('error-message');
    
    if (!username) {
        errorMessage.textContent = 'Please enter a username.';
        return;
    }

    try {
        const response = await fetch(`${API_URL}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            // Login successful
            localStorage.setItem('token', data.token);
            window.location.href = '/public/dashboard.html';
        } else {
            // --- THIS IS THE CORRECTED LOGIC ---
            // Check for our special redirect action from the backend
            if (data.action === 'redirect_to_signup') {
                window.location.href = `/public/signup.html?username=${username}`;
            } else {
                errorMessage.textContent = data.error;
            }
        }
    } catch (error) {
        errorMessage.textContent = 'An error occurred. Please try again.';
    }
}

function logout() {
    localStorage.removeItem('token');
    window.location.href = '/';
}

async function loadDashboard() {
    await fetchUsers();
    await fetchAppointments();
    await loadProfile(); // <-- FIX: Ensures the profile form is filled on load
}

async function fetchUsers() {
    const userListDiv = document.getElementById('user-list');
    const token = localStorage.getItem('token');
    const currentUser = parseJwt(token); // Get the current logged-in user

    const response = await fetch(`${API_URL}/api/users`, {
        headers: { 'Authorization': `Bearer ${token}` }
    });
    const users = await response.json();
    
    // FIX: Filter out the current user from the list
    userListDiv.innerHTML = users
        .filter(user => user.username !== currentUser.usr)
        .map(user =>
            `<input type="checkbox" name="participants" value="${user.id}"> ${user.name} (${user.username})<br>`
        ).join('');
}

async function fetchAppointments() {
    const appointmentList = document.getElementById('appointment-list');
    const token = localStorage.getItem('token');
    const response = await fetch(`${API_URL}/api/appointments`, {
        headers: { 'Authorization': `Bearer ${token}` }
    });
    const appointments = await response.json();
    if (appointments && appointments.length > 0) {
        appointmentList.innerHTML = appointments.map(appt =>
            `<li><b>${appt.title}</b>: Starts at ${appt.start_time_local}</li>`
        ).join('');
    } else {
        appointmentList.innerHTML = '<li>No upcoming appointments.</li>';
    }
}

async function createAppointment() {
    const title = document.getElementById('title').value;
    const localStartTime = document.getElementById('start_time').value;
    const startTimeUTC = new Date(localStartTime).toISOString();
    
    const duration = parseInt(document.getElementById('duration').value);
    const checkedBoxes = document.querySelectorAll('input[name="participants"]:checked');
    const participant_ids = Array.from(checkedBoxes).map(cb => cb.value);

    const formMessage = document.getElementById('form-message');
    const token = localStorage.getItem('token');

    try {
        const response = await fetch(`${API_URL}/api/appointments`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                title: title,
                start_time: startTimeUTC,
                duration_minutes: duration,
                participant_ids: participant_ids
            })
        });
        const data = await response.json();
        if (response.ok) {
            formMessage.style.color = 'green';
            formMessage.textContent = data.message;
            fetchAppointments();
        } else {
            formMessage.style.color = 'red';
            formMessage.textContent = `Error: ${data.error}`;
        }
    } catch (error) {
        formMessage.style.color = 'red';
        formMessage.textContent = 'A network error occurred.';
    }
}

async function loadProfile() {
    const token = localStorage.getItem('token');
    try {
        const response = await fetch(`${API_URL}/api/profile`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (response.ok) {
            const user = await response.json();
            document.getElementById('profile-name').value = user.name;
            document.getElementById('profile-timezone').value = user.preferred_timezone;
        }
    } catch (error) {
        console.error("Failed to load profile:", error);
    }
}

async function updateProfile() {
    const token = localStorage.getItem('token');
    const name = document.getElementById('profile-name').value;
    const preferred_timezone = document.getElementById('profile-timezone').value;
    const profileMessage = document.getElementById('profile-message');

    try {
        const response = await fetch(`${API_URL}/api/profile`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ name, preferred_timezone })
        });

        const data = await response.json();
        if (response.ok) {
            profileMessage.style.color = 'green';
            profileMessage.textContent = data.message;
            localStorage.setItem('token', data.token);
        } else {
            profileMessage.style.color = 'red';
            profileMessage.textContent = `Error: ${data.error}`;
        }
    } catch (error) {
        profileMessage.style.color = 'red';
        profileMessage.textContent = 'A network error occurred.';
    }
}