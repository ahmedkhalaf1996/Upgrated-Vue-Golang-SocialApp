describe("Auth Component Tests", () => {
    beforeEach(() => {
        cy.visit('/Auth')
    })

    describe('Page Layout', () => {
        it('should display signin and signup cards', () => {
            // Check for cards exists (should work on any viewport)
            cy.get('.q-card').should('have.length.at.least', 1)

            // Check titles - these exist in both desktop and mobile versions
            cy.contains('Sign In').should('exist')
            // On mobile, it shows in tabs, on desktop in card header
        })

        it('should have proper layout structure', () => {
            cy.get('.q-page').should('exist')
            cy.get('.constrain').should('exist')
        })
    })

    describe('Signin Form', () => {
        it('should display signin form elements', () => {
            // Set desktop viewport to ensure we see desktop layout
            cy.viewport(1200, 800)
            
            cy.get('.col-5 input').should('have.length', 2)
            cy.get('.col-5').contains('Sign In').should('be.visible')
            cy.get('.col-5 button[type="submit"]').should('be.visible')
        })

        it('should allow typing in signin inputs', () => {
            cy.viewport(1200, 800)
            
            cy.get('.col-5 input').eq(0)
                .type('test@example.com')
                .should('have.value', 'test@example.com')
            
            cy.get('.col-5 input').eq(1)
              .type('password123')
              .should('have.value', 'password123')
        })

        it('should have password input type', () => {
            cy.viewport(1200, 800)
            
            cy.get('.col-5 input').eq(1)
              .should('have.attr', 'type', 'password')
        })
    })

    describe('Signup Form', () => {
        it('should display all signup form elements', () => {
            cy.viewport(1200, 800)
            
            // Check all signup inputs exist
            cy.get('.col-7 input').should('have.length', 4)
            cy.contains('Your First Name *').should('be.visible')
            cy.contains('Your Last Name *').should('be.visible')
            cy.contains('Your Email *').should('be.visible')
            cy.contains('Your Password *').should('be.visible')

            // Check signup button
            cy.contains('Create New Account').should('be.visible')
        })

        it('should allow typing in all signup inputs', () => {
            cy.viewport(1200, 800)
            
            cy.get('.col-7 input').eq(0)
              .type('John')
              .should('have.value', 'John')

            cy.get('.col-7 input').eq(1)
              .type('Doe')
              .should('have.value', 'Doe')

            cy.get('.col-7 input').eq(2)
              .type('john@example.com')
              .should('have.value', 'john@example.com')

            cy.get('.col-7 input').eq(3)
              .type('password123')
              .should('have.value', 'password123')
        })

        it('should have correct button colors', () => {
            cy.viewport(1200, 800)
            
            // signin button
            cy.get('.col-5 .q-btn').should('have.class', 'bg-primary')

            // signup button  
            cy.get('.col-7 .q-btn').should('have.class', 'bg-positive')
        })
    })

    describe('Form Interactions', () => {
        it('should handle empty input submissions', () => {
            cy.viewport(1200, 800)
            
            // Test signin form validation
            cy.get('.col-5 button[type="submit"]').click()
            cy.get('.q-notification').should('contain.text', 'Email is required')

            // Wait for notification to clear
            cy.wait(2000)

            // Test signup form validation - just check that validation happens
            cy.get('.col-7 button[type="submit"]').click()
            cy.get('.q-notification').should('be.visible')
            cy.get('.q-notification').should('contain.text', 'is required')
        })
    })

    describe('Mobile Layout', () => {
        it('should display mobile elements on small screens', () => {
            cy.viewport(375, 667)
            
            // Mobile tabs should be visible
            cy.get('.mobile-tabs').should('be.visible')
            
            // Desktop layout should be hidden
            cy.get('.lt-sm-hide').should('not.be.visible')
            
            // Mobile layout should be visible
            cy.get('.mobile-layout').should('be.visible')
        })

        it('should switch between forms on mobile using tabs', () => {
            cy.viewport(375, 667)
            
            // Click on tab elements (not by name attribute, but by visible text)
            cy.get('.mobile-tabs').contains('Sign In').should('be.visible')
            cy.get('.mobile-tabs').contains('Sign Up').should('be.visible')
            
            // Click signup tab
            cy.get('.mobile-tabs').contains('Sign Up').click()
            
            // Verify we can see signup content
            cy.contains('Join Us').should('be.visible')
            
            // Click signin tab  
            cy.get('.mobile-tabs').contains('Sign In').click()
            
            // Verify we can see signin content
            cy.contains('Welcome Back').should('be.visible')
        })

        it('should allow form input on mobile', () => {
            cy.viewport(375, 667)
            
            // Test signin form on mobile
            cy.get('.mobile-tabs').contains('Sign In').click()
            
            // Find and fill mobile signin inputs
            cy.get('.mobile-card input').first()
                .type('test@example.com')
                .should('have.value', 'test@example.com')
        })
    })

    describe('Responsive Design', () => {
        it('should maintain layout on different screen sizes', () => {
            // Mobile view
            cy.viewport(375, 667)
            cy.get('.mobile-tabs').should('be.visible')

            // Tablet/Desktop view
            cy.viewport(768, 1024)
            cy.get('.col-5').should('be.visible')
            cy.get('.col-7').should('be.visible')

            // Large Desktop view
            cy.viewport(1200, 800)
            cy.get('.row').should('be.visible')
            cy.get('.col-5').should('be.visible')
            cy.get('.col-7').should('be.visible')
        })
    })
})