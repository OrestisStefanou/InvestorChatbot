# TODO:
1. Handle telegram message text limit
   2. splitting by new line is not good enough because it spams the user with messages
   3. We could try by new line but also check message size and keep adding until is close to the limit
2. Current message format problems
   3. The service could return markdown tables as well that we can't handle
   4. Transform message function could become overcomplicated.