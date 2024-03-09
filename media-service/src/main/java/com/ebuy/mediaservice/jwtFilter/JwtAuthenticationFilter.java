package com.ebuy.mediaservice.jwtFilter;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.Collection;
import java.util.List;

import javax.crypto.SecretKey;

import org.springframework.lang.NonNull;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.context.SecurityContext;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.web.authentication.WebAuthenticationDetailsSource;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.Cookie;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jws;
import io.jsonwebtoken.JwtException;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.security.Keys;

@Component
public class JwtAuthenticationFilter extends OncePerRequestFilter {

    @Override
    protected void doFilterInternal(@NonNull HttpServletRequest request,
            @NonNull HttpServletResponse response, @NonNull FilterChain filterChain)
            throws ServletException, IOException {
        Cookie[] cookies = request.getCookies();
        if (cookies == null) {
            filterChain.doFilter(request, response);
            return;
        }

        String jwtToken = null;
        for (Cookie cookie : cookies) {
            if (cookie.getName().equals("jwt")) {
                jwtToken = cookie.getValue();
                break;
            }
        }

        if (jwtToken == null || jwtToken.isEmpty() || SecurityContextHolder.getContext().getAuthentication() != null) {
            filterChain.doFilter(request, response);
            return;
        }

        Jws<Claims> claims;
        try {
            claims = verifyAndDecodeJwt(jwtToken);
        } catch (JwtException e) {
            System.out.println(e);
            filterChain.doFilter(request, response);
            return;
        }

        Integer userId = (Integer) claims.getPayload().get("user_id");
        User userDetails = new User(userId);
        if (userId > 0) {
            SecurityContext context = SecurityContextHolder.createEmptyContext();
            UsernamePasswordAuthenticationToken authToken = new UsernamePasswordAuthenticationToken(
                    userDetails, null, getGrantedAuthorities());
            authToken.setDetails(new WebAuthenticationDetailsSource().buildDetails(request));
            context.setAuthentication(authToken);
            SecurityContextHolder.setContext(context);
        }

        filterChain.doFilter(request, response);
    }

    private static Collection<? extends GrantedAuthority> getGrantedAuthorities() {
        List<GrantedAuthority> authorities = new ArrayList<>();

        // Add roles to the user
        authorities.add(new SimpleGrantedAuthority("ROLE_USER"));
        authorities.add(new SimpleGrantedAuthority("ROLE_ADMIN"));

        return authorities;
    }

    public Jws<Claims> verifyAndDecodeJwt(String jwtString) {
        Jws<Claims> jws;
        try {
            // String secretKeyString =
            // "qwertyuiopasdfghjklzxcvbnm123456qwertyuiopasdfghjklzxcvbnm123456";
            String secretKeyString = System.getenv("jwtSecret");
            SecretKey secretKey = Keys.hmacShaKeyFor(secretKeyString.getBytes(StandardCharsets.UTF_8));
            // Jwts.SIG.HS512.key(SecretKey);
            jws = Jwts.parser() // (1)
                    .verifyWith(secretKey)
                    .build()
                    .parseSignedClaims(jwtString);

        } catch (JwtException ex) { // (5)
            throw ex;
        }
        return jws;
    }
}