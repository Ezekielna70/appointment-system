# Technical Questions

### 1. Timezone Conflicts: How would you handle timezone conflicts between participants in an appointment?

Handling timezone conflicts is absolutely crucial for a scheduling application. My implementation addresses this with a foundational rule and could be expanded with more advanced, user-friendly features.

* **Current Implementation (Working Hours Validation):** The system's core logic for handling conflicts is to enforce "working hours." When an appointment is created with invitees, the system converts the proposed UTC time to each invitee's local timezone (based on their profile) and checks if it falls within a reasonable working window (08:00 - 17:00). If the proposed time is outside of these hours for **any** invitee, the API rejects the request and informs the creator of the conflict. This prevents users from scheduling meetings for others at inconvenient times.

* **Advanced Strategy 1 (Visual Feedback):** For a real product, I would enhance the frontend UI. As the organizer selects a time slot, the interface would instantly show the corresponding local time for each participant. So that the user would be easier to visualize the time alocation for other user



---

### 2. Database Optimization: How can you optimize database queries to efficiently fetch user-specific appointments?

To make sure the app stays fast even with thousands of users and appointments, we use two main tricks to keep the database super fast.

#### **Using an Index (Like a Book's Index)**

Imagine trying to find a topic in a huge book without an index at the back. You'd have to flip through every single page. A database index does the exact same thing for your data, letting you jump straight to the information you need.

* **For Users:** We put an index on the `username`. When you try to log in, the database uses this index to instantly find your username instead of searching through every user. This makes logging in very fast and also prevents two people from having the same username.

* **For Appointments:** We use a smart, two-part index on `participants` and `start_time`. When you ask for your upcoming appointments, the database first uses the index to instantly grab *only* the appointments you're a part of. Then, within that much smaller group, it quickly checks the date to find the upcoming ones. This is much faster than looking through every appointment in the system.

#### **Only Asking for What You Need (Projection)**

Imagine you're ordering at a restaurant and you just want a burger. You wouldn't ask the waiter to bring every single item from the kitchen to your table. You'd just say, "I'll have the burger, please."

This is the same idea. When the dashboard shows the "Upcoming Appointments" list, it only needs the **title** and the **time**. So, instead of asking the database for the *entire* appointment record, we tell it: "**Just give us the title and the start time.**" This means less data has to travel from the database to our app, which makes the page load faster.

---

### 3. Additional Features: If this application were to become a real product, what additional features would you implement? Why?

To ensure this application into a working and competitive product, I would implement the following features:

* **Notifications and Reminders:**
    * **What:** Automatically send email or in-app notifications to users when they are invited to an appointment. 
    * **Why:** This is a very basic feature for any scheduling tool. It ensures users are aware of their schedule and reduces the chance of missed appointments.

* **Recurring Appointments:**
    * **What:** Allow users to create appointments that can repeat on a schedule (e.g., daily, weekly on a specific day, monthly) by giving the choice up to the user's hand.
    * **Why:** Many professional meetings (like weekly team syncs or monthly reviews) are recurring. This feature is definitely a quality-of-life upgrade that can save users the manual effort of creating the same event over and over.


---

### 4. Session Management: How would you manage user sessions securely while keeping them lightweight?

My current implementation uses JSON Web Tokens (JWTs), which are inherently stateless and lightweight. To manage them securely, I would follow these best practices:

* **Use Strong Secret Keys**: The key used to sign the tokens is stored securely as an environment variable and is never exposed in the code.

* **Limit The Login Attempts**: By limiting the login attempts, it prevents a bruteforce way of login and prevents malicious attempts by unkwown user.

* **Secure Storage**: The JWT should be stored in an `HttpOnly`, `SameSite` cookie on the client-side. This is more secure than `localStorage` because it prevents the token from being stolen via Cross-Site Scripting (XSS) attacks.