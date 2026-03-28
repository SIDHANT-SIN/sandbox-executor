
<h1>Sandbox Executor</h1>
<p><strong>Repo:</strong> sandbox-executor | Built in Go | Multi-language code execution (Python, C++, Java)</p>

<h2>Overview</h2>
<p>Sandbox Executor is a secure and efficient backend service for executing code in multiple languages. It runs Python, C++, and Java programs inside isolated Docker containers to ensure system safety and concurrency control.</p>

<h2>Features</h2>

<h3>Multi-language Execution</h3>
<ul>
    <li>Python </li>
    <li>C++ </li>
    <li>Java </li>
</ul>

<h3>Docker Sandbox</h3>
<ul>
    <li>Memory limited to 256MB</li>
    <li>CPU limited to 0.5 cores</li>
    <li>Time Limit 2 sec</li>
    <li>Network disabled</li>
    <li>Process limit 64</li>
</ul>
<p>Prevents fork bombs, memory abuse, and unauthorized network usage.</p>

<h3>Compile vs Run Separation</h3>
<ul>
    <li>Compile step is separate and not timed</li>
    <li>Run step is timed (here 2s timeout)</li>
    <li>Matches behavior of real online judges</li>
</ul>

<h3>Error Classification</h3>
<table border="1" cellpadding="5">
    <tr>
        <th>Case</th>
        <th>Status</th>
    </tr>
    <tr>
        <td>Syntax error</td>
        <td>compile_error</td>
    </tr>
    <tr>
        <td>Runtime crash</td>
        <td>runtime_error</td>
    </tr>
    <tr>
        <td>Infinite loop</td>
        <td>timeout</td>
    </tr>
    <tr>
        <td>Compile hang</td>
        <td>compile_timeout</td>
    </tr>
    <tr>
        <td>Infrastructure issue</td>
        <td>server_error</td>
    </tr>
</table>

<h3>Concurrency Control</h3>
<ul>
    <li>Semaphore limits max concurrent executions</li>
    <li>Prevents CPU/memory exhaustion and Docker crashes</li>
    <li>Rate limiting per client implemented with exponential backoff + jitter</li>
</ul>

<h3>Security</h3>
<ul>
    <li>Authorization via secret key in header (`Authorization: Bearer <SECRET_KEY>`) </li>
    <li>Rate limiting prevents abuse and throttles excessive requests </li>
</ul>

<h3>Request Isolation</h3>
<ul>
    <li>Each request runs in a temporary directory</li>
    <li>Prevents file clashes between executions</li>
    <li>Docker containers cleaned up automatically with <code>--rm</code></li>
</ul>

<h2>Getting Started</h2>
<ol>
    <li>Clone the repo</li>
    <li>Create a <code>.env</code> file with required secrets (look into .env.example)</li>
    <li>Run the server:
        <pre><code>go run main.go</code></pre>
    </li>
    <li>Call the <code>/execute</code> endpoint with code and language parameters, including the authorization header</li>
</ol>

<h2>Example Request</h2>
<pre><code>POST /execute
Authorization: Bearer YOUR_SECRET_KEY
Content-Type: application/json

{
  "language": "python",
  "code": "print('Hello World')"
}
</code></pre>

<h2>Notes</h2>
<ul>
    <li>Designed for safe execution of untrusted code</li>
    <li>Timeouts, sandboxing, and concurrency control make it production-ready</li>
    <li>Extensible to add more languages in the future</li>
</ul>

<h2>License</h2>
<p>MIT License</p>
