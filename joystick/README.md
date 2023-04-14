# Joystick module

Handles joystick input.

Unlike other modules, this module relies on glfw to handle joystick input. This means that you must use glfw and call `glfwPollEvents()` in your main loop for joystick input to work.