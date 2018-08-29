%% TASK 5.1a -- Estimate PSD %
%figNum = 1;             % Figure number-counter
load('wave.mat');       % Load wave disturbance

F_s = 10;
window = 4096;
noverlap = [];
nfft = [];  
[S_psi,f] = pwelch(psi_w(2,:).*(pi/180),window,noverlap,nfft,F_s);
omega = 2*pi.*f;
S_psi = S_psi./(2*pi);


%% TASK 5.2c ---- Find omega_0
% Plot estimated PSD
width = 20; % cm
height = 10; % cm
fontsize = 10; %points

set(0,'DefaultTextInterpreter', 'latex') %Interpret (most) text as latex. 

plot(omega,S_psi, 'LineWidth', 2)
axis([0 2 -0.00005 16*10^(-4)])
hold on
xlabel('$\omega$ [$\frac{rad}{s}$]')
ylabel('$S_{\psi_{w}}(\omega)$ [rad]')
title(['Estimated power spectral density fuction of $S_{\psi_{w}}(\omega)$ '...
   ])
grid on; 

ax = gca; ax.XTick = [0:pi/8:2]; %format the x axis to pi
ax.XTickLabel = {'$0$', '$\frac{\pi}{8}$', '$\frac{\pi}{4}$', ...
    '$\frac{3\pi}{8}$', '$\frac{\pi}{2}$','$\frac{5\pi}{8}$', ...
    '$\frac{3\pi}{4}$'};
ax.TickLabelInterpreter = 'latex';



%% Finding max 5.2b

%Make some comments here
[maxPSD,  frequency_index ] = max( S_psi )
omega_0 = omega( frequency_index ) 




%% TASK 5.2.d ---- Finding lambda
sigma = sqrt(maxPSD);

% lsqcurvefit (least-squares nonlinear curve fitting)


  P_psi_fun = @(lambda,omega) ...
       (4*lambda^2*omega_0^2*sigma^2*omega.^2) ./ ...
 (omega.^4 + (2*lambda^2 - 1)*2*omega_0^2*omega.^2 + ...
omega_0^4);

lambda0 = 10;
lb=0;
ub=10;
lambda = lsqcurvefit(P_psi_fun,lambda0,omega,S_psi,lb,ub);
K_w = 2*lambda*omega_0*maxPSD;
P_psi = P_psi_fun(lambda,omega);

 %%% Comparison plot of estimate and analytical
figure
plot(omega, P_psi, 'r')
hold on
plot(omega, S_psi, 'b')
legend('P_{\psi_w}','S_{\psi_w}')
axis([0 2 -0.00005 16*10^(-4)])
hold on
xlabel('t[s]')
ylabel('$S_{\psi_{w}}(\omega)$ [rad]')
title(['Estimated power spectral density fuction of $S_{\psi_{w}}(\omega)$ '...
   ])
grid on; 

ax = gca; ax.XTick = [0:pi/8:2]; %format the x axis to pi
ax.XTickLabel = {'$0$', '$\frac{\pi}{8}$', '$\frac{\pi}{4}$', ...
    '$\frac{3\pi}{8}$', '$\frac{\pi}{2}$','$\frac{5\pi}{8}$', ...
    '$\frac{3\pi}{4}$'};
ax.TickLabelInterpreter = 'latex';