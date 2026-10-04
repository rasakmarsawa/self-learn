type NavbarProps = {
  username: string
  onLogout: () => void
}

function Navbar({ username, onLogout }: NavbarProps) {
  return (
    <header className="navbar">
      <span className="navbar-brand">Rabbit MQ</span>

      <nav className="navbar-nav" aria-label="Main">
        <a className="navbar-link" href="#dashboard">
          Dashboard
        </a>
        <a className="navbar-link" href="#messages">
          Messages
        </a>
        <a className="navbar-link" href="#workers">
          Workers
        </a>
        <a className="navbar-link" href="#settings">
          Settings
        </a>
      </nav>

      <div className="navbar-user">
        <span className="navbar-username">{username}</span>
        <button type="button" className="navbar-logout" onClick={onLogout}>
          Logout
        </button>
      </div>
    </header>
  )
}

export default Navbar
